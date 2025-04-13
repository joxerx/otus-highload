package main

import (
	"log"
	"os"

	"net/http"

	"otus-highload-counter/redis"
	"otus-highload-counter/router"
)

func main() {
	redis.InitRedis()
	redis.CreateEventGroup()
	go redis.StartEventConsumer()
	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GO_PORT"), r))
}
