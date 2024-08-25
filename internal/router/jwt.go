package router

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zimmah/chirpy/internal/database"
)

type Claims struct {
	jwt.RegisteredClaims
}

type JWT struct {
	Password 			string `json:"password"`
	Email 				string `json:"email"`
	ExpiresInSeconds 	*int `json:"expires_in_seconds"`
}

func (cfg *apiConfig) generateJWT(user database.User, expiresInSeconds *int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	if expiresInSeconds != nil {
		customExpiration := time.Duration(*expiresInSeconds) * time.Second
		if customExpiration > 24*time.Hour {
			customExpiration = 24 * time.Hour
		}
		expirationTime = time.Now().Add(customExpiration)
	}

	claims := Claims {
		jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(expirationTime.UTC()),
		Subject: fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(cfg.jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (cfg *apiConfig) validateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return cfg.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid && !claims.ExpiresAt.Time.UTC().Before(time.Now().UTC()) {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}