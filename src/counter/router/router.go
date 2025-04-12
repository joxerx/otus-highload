package router

import (
	"otus-highload-counter/handlers/counter"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/counter/list", counter.ListCounterHandler).Methods("GET")
	router.HandleFunc("/counter/{dialogId}", counter.DialogCounterHandler).Methods("GET")

	return router
}
