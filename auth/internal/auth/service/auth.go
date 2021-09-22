package service

import (
	"errors"
	"log"
	"math/rand"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ulwisibaq/efishery/auth/internal/auth"
	"github.com/ulwisibaq/efishery/auth/internal/models"
	"github.com/ulwisibaq/efishery/auth/pkg/constants"
)

type AuthService struct {
	authRepository auth.MysqlRepository
}

func NewAuthService(
	authRepository auth.MysqlRepository,
) auth.Service {
	return &AuthService{
		authRepository: authRepository,
	}
}

func (as AuthService) Register(userReq models.UserRequest) (resp models.User, err error) {

	password := as.generatePassword()

	roleId, err := as.authRepository.GetRoleIdByName(userReq.Role)
	if err != nil {
		return
	}
	if roleId == 0 {
		err = errors.New(constants.ErrRoleNotFound)
		return
	}

	userReq.Created = time.Now()
	userReq.Password = password
	userReq.RoleId = roleId

	id, err := as.authRepository.CreateUser(userReq)
	if err != nil {
		return
	}

	resp.ID = id
	resp.Name = userReq.Name
	resp.Phone = userReq.Phone
	resp.Role = userReq.Role
	resp.Password = userReq.Password
	resp.Created = userReq.Created

	return

}

func (as AuthService) Login(login models.LoginRequest) (resp models.LoginResponse, err error) {

	user, err := as.authRepository.GetUserByPhone(login.Phone)
	if err != nil {
		return
	}
	if user.ID == 0 {
		err = errors.New(constants.ErrUserNotFound)
		return
	}

	if login.Password != user.Password {
		err = errors.New(constants.ErrIncorrectPassword)
		return
	}

	token := as.generateToken(user)

	resp.Token = token

	return
}

func (as AuthService) Verify(token string) (*models.UserClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			err := errors.New(constants.ErrInvalidToken)
			return nil, err
		}
		return []byte(constants.JwtSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &models.UserClaims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		errExp := errors.New(constants.ErrExpiredToken)
		if ok && errors.Is(verr.Inner, errExp) {
			return nil, errExp
		}
		return nil, errors.New(constants.ErrInvalidToken)
	}

	payload, ok := jwtToken.Claims.(*models.UserClaims)
	if !ok {
		return nil, errors.New(constants.ErrInvalidToken)
	}

	return payload, nil
}

func (as AuthService) generatePassword() string {
	byteString := make([]byte, 4)
	for i := range byteString {
		letterBytes := constants.LetterBytes
		byteString[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(byteString)
}

func (as AuthService) generateToken(user models.User) (userToken string) {
	JWTSigningMethod := jwt.SigningMethodHS256
	JWTSignatureKey := []byte(constants.JwtSecret)

	expireTime := time.Now().Add(time.Duration(30) * time.Minute)

	claims := models.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    constants.JwtIssuer,
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Name:    user.Name,
		Phone:   user.Phone,
		Role:    user.Role,
		Created: user.Created,
	}

	token := jwt.NewWithClaims(
		JWTSigningMethod,
		claims,
	)

	signedToken, err := token.SignedString(JWTSignatureKey)
	if err != nil {
		log.Fatal(err)
		return
	}

	return signedToken
}
