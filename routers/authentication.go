package routers

import (
	"api.jwt.auth/controllers"
	"api.jwt.auth/core/authentication"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func SetAuthenticationRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/token-auth", controllers.Login).Methods("POST")
	router.Handle("/refresh-token-auth", negroni.New(negroni.HandlerFunc(authentication.RequireTokenAuthentication), negroni.HandlerFunc(controllers.RefresfhToken))).Methods("GET")
	return router
}
