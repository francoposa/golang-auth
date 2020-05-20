package resources

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}
