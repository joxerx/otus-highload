package main

import (
	"log"
	"os"

	"net/http"

	"otus-highload-messages/redis"
	"otus-highload-messages/router"
)

func main() {
	// redis.InitRedis()
	// db.InitDB()
	// defer func() {
	// 	if db.MasterDB != nil {
	// 		db.MasterDB.Close()
	// 	}
	// 	if db.BalancerDB != nil {
	// 		db.BalancerDB.Close()
	// 	}
	// 	log.Println("Database connections closed.")
	// }()

	redis.InitRedis()
	r := router.NewRouter()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GO_PORT"), r))
}
