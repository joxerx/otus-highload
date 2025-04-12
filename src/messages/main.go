package main

import (
	"log"
	"os"

	"net/http"

	"otus-highload-messages/redis"
	"otus-highload-messages/router"
)

func main() {
	redis.InitRedis()
	go redis.StartEventConsumer()

	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GO_PORT"), r))
}
