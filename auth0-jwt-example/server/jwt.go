package server

import (
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("secret")

func NewJWTMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(
		jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return mySigningKey, nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		},
	)
}

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to hold our token claims
	claims := token.Claims.(jwt.MapClaims)

	// Set token claims
	claims["admin"] = true
	claims["name"] = "Franco Posa"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Sign the token with our secret
	tokenString, _ := token.SignedString(mySigningKey)

	// Finally, write the token to the response
	w.Write([]byte(tokenString))
}
