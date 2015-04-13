package services

import (
	"api.jwt.auth/api/parameters"
	"api.jwt.auth/core/authentication"
	"api.jwt.auth/services/models"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

func Login(requestUser *models.User) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()

	if authBackend.Authenticate(requestUser) {
		token := parameters.TokenAuthentication{authBackend.GenerateToken(requestUser)}
		response, _ := json.Marshal(token)
		return http.StatusOK, response
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(requestUser *models.User) []byte {
	authBackend := authentication.InitJWTAuthenticationBackend()
	token := parameters.TokenAuthentication{authBackend.GenerateToken(requestUser)}
	response, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}
	return response
}

func Logout(req *http.Request) error {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenRequest, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authBackend.Logout(tokenString, tokenRequest)
}
