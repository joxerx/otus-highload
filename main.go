package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Birthdate string `json:"birthdate"`
	Biography string `json:"biography"`
	City      string `json:"city"`
	Password  string `json:"password"`
}

func main() {
	db = initDB()
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/user/register", registerUserHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("GO_PORT"), router))
}

func initDB() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
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

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	log.Println("Db connected!")
	return db
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	birthdate, err := time.Parse("2006-01-02", newUser.Birthdate)
	if err != nil {
		http.Error(w, "Invalid birthdate format. Use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(
		"INSERT INTO users (first_name, last_name, birthdate, biography, city, password) VALUES ($1, $2, $3, $4, $5, $6)",
		newUser.FirstName, newUser.LastName, birthdate, newUser.Biography, newUser.City, hashedPassword,
	)
	if err != nil {
		http.Error(w, "Error inserting user into db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
