package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ulwisibaq/efishery/auth/pkg/constants"
	"github.com/ulwisibaq/efishery/commodity/internal/commodity"
)

type commodityHandler struct {
	commodityService commodity.Service
}

func NewCommodityHandler(
	commodityService commodity.Service,
) *commodityHandler {
	return &commodityHandler{
		commodityService: commodityService,
	}
}

func (ch *commodityHandler) GetCommodity(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")

	splitToken := strings.Split(authToken, "Bearer")
	if len(splitToken) != 2 {
		respondWithError(w, http.StatusBadRequest, constants.ErrInvalidAuthorization)
		return
	}

	token := strings.TrimSpace(splitToken[1])

	_, err := ch.commodityService.Verify(token)
	if err != nil {
		if err.Error() == constants.ErrExpiredToken {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if err.Error() == constants.ErrInvalidToken {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
	}

	commodities, err := ch.commodityService.GetCommodity()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, commodities)
}

func (ch *commodityHandler) Verify(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")

	splitToken := strings.Split(authToken, "Bearer")
	if len(splitToken) != 2 {
		respondWithError(w, http.StatusBadRequest, constants.ErrInvalidAuthorization)
		return
	}

	token := strings.TrimSpace(splitToken[1])

	resp, err := ch.commodityService.Verify(token)
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
