package router

import (
	"otus-highload/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/user/register", handlers.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/user/get/{id}", handlers.GetUserHandler).Methods("GET")
	router.HandleFunc("/user/login", handlers.LoginHandler).Methods("POST")
	return router
}
