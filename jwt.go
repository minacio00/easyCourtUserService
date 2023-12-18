package main

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// secret key
var Jwtkey = []byte(viper.GetString("JWT_SECRET_KEY"))

func generateJWT(username string) (string, error) {
	// Create a new set of claims
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create the token using the claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(Jwtkey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
