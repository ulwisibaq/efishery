package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	jwt.StandardClaims
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Role    string    `json:"role"`
	Created time.Time `json:"created_at"`
}
