package router

import (
	"fmt"
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zimmah/chirpy/internal/database"
)

type JWT struct {
	Password 			string `json:"password"`
	Email 				string `json:"email"`
	ExpiresInSeconds 	int `json:"expires_in_seconds"`
}

func (cfg *apiConfig) generateJWT(user database.User, expiresInSeconds int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	if expiresInSeconds != 0 {
		customExpiration := time.Duration(expiresInSeconds) * time.Second
		if customExpiration > 24*time.Hour {
			customExpiration = 24 * time.Hour
		}
		expirationTime = time.Now().Add(customExpiration)
	}

	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(expirationTime.UTC()),
		Subject: fmt.Sprintf("%d", user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(cfg.jwtSecret)
}

func (cfg *apiConfig) validateJWT(tokenString string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (any, error) { return []byte(cfg.jwtSecret), nil },
	)
	if err != nil { return "", err }

	userIDString, err := token.Claims.GetSubject()
	if err != nil { return "", err }

	issuer, err := token.Claims.GetIssuer()
	if err != nil { return "", err }
	if issuer != string("chirpy") { return "", errors.New("invalid issuer") }

	return userIDString, nil
}