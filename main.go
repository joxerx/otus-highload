package main

import (
	"log"
	"net/http"
	"os"

	"otus-highload/db"
	"otus-highload/router"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GO_PORT"), r))
}
