package router

import (
	"otus-highload-messages/handlers/dialog"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/dialog/{userId}/send", dialog.SendDialogHandler).Methods("POST")
	router.HandleFunc("/dialog/{userId}/list", dialog.ListDialogHandler).Methods("GET")

	return router
}
