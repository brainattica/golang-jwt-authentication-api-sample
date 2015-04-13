package controllers

import (
	"api.jwt.auth/core/authentication"
	"api.jwt.auth/services"
	"api.jwt.auth/services/models"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := services.Login(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func RefresfhToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(services.RefreshToken())
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
