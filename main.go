package main

import (
	"log"
	"os"
	"syscall"

	"net/http"
	"os/signal"

	"otus-highload/db"
	"otus-highload/redis"
	"otus-highload/router"
	"otus-highload/utils"
)

func main() {
	redis.InitRedis()
	db.InitDB()
	defer func() {
		if db.MasterDB != nil {
			db.MasterDB.Close()
		}
		for _, slaveDB := range db.SlaveDBs {
			if slaveDB != nil {
				slaveDB.Close()
			}
		}
		log.Println("Database connections closed.")
	}()

	stopChan := make(chan struct{})

	go redis.StartTaskConsumer(stopChan)
	go utils.StartPeriodicTask()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received shutdown signal. Cleaning up...")
		close(stopChan)
	}()

	// log.Println("Warming up the cache... Waiting 5 minutes.")
	// time.Sleep(5 * time.Minute)
	// log.Println("Ready to serve requests.")

	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GO_PORT"), r))
}
