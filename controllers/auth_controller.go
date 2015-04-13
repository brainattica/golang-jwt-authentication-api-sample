package controllers

import (
	"api.jwt.auth/api/parameters"
	"api.jwt.auth/core/authentication"
	"api.jwt.auth/services/models"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	request_user := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request_user)

	authBackend := authentication.InitJWTAuthenticationBackend()

	if authBackend.Authenticate(request_user) {
		token := parameters.TokenAuthentication{authBackend.GenerateToken()}
		response, _ := json.Marshal(token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}
}

func RefresfhToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	token := parameters.TokenAuthentication{authBackend.GenerateToken()}
	response, _ := json.Marshal(token)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	tokenString := r.Header.Get("Authorization")

	err = authBackend.Logout(tokenString, token)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(http.StatusInternalServerError)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		fmt.Println(http.StatusOK)
		w.WriteHeader(http.StatusOK)
	}
}
