package main

import (
	"log"
	"os"

	"net/http"

	"otus-highload-messages/metrics"
	"otus-highload-messages/redis"
	"otus-highload-messages/router"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metrics.Init()
	redis.InitRedis()
	redis.CreateEventGroup()
	go redis.StartEventConsumer()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GO_PORT"), r))
}
