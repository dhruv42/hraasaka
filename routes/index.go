package routes

import (
	"net/http"

	"github.com/dhruv42/hraasaka/controllers"
	"github.com/gorilla/mux"
)

func New() http.Handler {
	route := mux.NewRouter()

	api := route.PathPrefix("/api").Subrouter()

	api.HandleFunc("/shorten", controllers.ShortenUrl).Methods("POST")
	api.HandleFunc("/redirect", controllers.RedirectUrl).Methods("POST")

	return route
}
