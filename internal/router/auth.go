package router

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zimmah/chirpy/internal/database"
)

type JWT struct {
	Password 			string `json:"password"`
	Email 				string `json:"email"`
}

func (cfg *apiConfig) generateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(maxTokenLifetime)

	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(expirationTime.UTC()),
		Subject: fmt.Sprintf("%d", userID),
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

func (cfg *apiConfig) generateRefreshToken(userID int) (string, error){
	randomData := make([]byte, 32)
	_, err := rand.Read(randomData)
	if err != nil {return "", err}

	token := hex.EncodeToString(randomData)

	expiresAt := time.Now().Add(maxRefreshTokenLifetime).Unix()
	err = database.DBPointer.UpdateUserToken(userID, int(expiresAt), token)
	if err != nil {
		return "", err
	}

	return token, nil
}