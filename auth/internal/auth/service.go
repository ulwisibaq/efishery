package auth

import "github.com/ulwisibaq/efishery/auth/internal/models"

type Service interface {
	Register(userReq models.UserRequest) (resp models.User, err error)
	Login(login models.LoginRequest) (resp models.LoginResponse, err error)
	Verify(token string) (*models.UserClaims, error)
}
