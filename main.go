package main

import (
	"log"
	"net/http"
	"os"

	"otus-highload/db"
	"otus-highload/redis"
	"otus-highload/router"
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

	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GO_PORT"), r))
}
