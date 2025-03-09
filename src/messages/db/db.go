package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var MasterDB *sql.DB
var SlaveDBs map[string]*sql.DB

func InitDB() {
	MasterDB = connectToDB(os.Getenv("POSTGRES_HOST"))

	slaveIDs := []string{"slave-1", "slave-2"} // Define slave IDs
	SlaveDBs = make(map[string]*sql.DB)

	for _, slaveID := range slaveIDs {
		SlaveDBs[slaveID] = connectToDB(fmt.Sprintf("db-%s", slaveID))
	}

	log.Println("Databases initialized!")
}

func connectToDB(host string) *sql.DB {
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Configure .env before starting service\nTip: call `make init`")
	}

	connectionStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	log.Printf("Connected to database %s", host)
	return db
}
