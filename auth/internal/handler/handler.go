package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ulwisibaq/efishery/auth/internal/auth"
	"github.com/ulwisibaq/efishery/auth/internal/models"
	"github.com/ulwisibaq/efishery/auth/pkg/constants"
)

type authHandler struct {
	authService auth.Service
}

func NewAuthHandler(
	authService auth.Service,
) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

func (ah *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userRequest := models.UserRequest{}

	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	resp, err := ah.authService.Register(userRequest)
	if err != nil {
		if err.Error() == constants.ErrRoleNotFound {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.Response{
		Data: resp,
	}

	respondWithJSON(w, http.StatusCreated, response)

}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginReq := models.LoginRequest{}

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	resp, err := ah.authService.Login(loginReq)
	if err != nil {
		if err.Error() == constants.ErrUserNotFound {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err.Error() == constants.ErrIncorrectPassword {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.Response{
		Data: resp,
	}

	respondWithJSON(w, http.StatusCreated, response)

}

func (ah *authHandler) Verify(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	splitToken := strings.Split(authToken, "Bearer")
	if len(splitToken) != 2 {
		respondWithError(w, http.StatusBadRequest, constants.ErrInvalidAuthorization)
		return
	}

	token := strings.TrimSpace(splitToken[1])

	resp, err := ah.authService.Verify(token)
	if err != nil {
		if err.Error() == constants.ErrExpiredToken {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if err.Error() == constants.ErrInvalidToken {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, resp)

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"Error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
