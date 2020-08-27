package utils

import (
	"api_new/logger"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Claims JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func init() {}

func GenerateJwtToken(email string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			Id:        email,
			Issuer:    "system", // contains provider for authenticate
			IssuedAt:  time.Now().UTC().Unix(),
			NotBefore: time.Now().UTC().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(24)).Unix(),
		},
	})
	token, err := jwtToken.SignedString([]byte("JWTsecret"))
	if err != nil {
		logger.Error("Generate token error: ", err)
		return "", err
	}
	logger.Info("token: ", token)
	return token, nil
}
