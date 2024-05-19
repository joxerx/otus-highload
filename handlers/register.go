package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"otus-highload/db"
	"otus-highload/models"
	"otus-highload/utils"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newUser.FirstName == "" || newUser.LastName == "" || newUser.Birthdate == "" || newUser.Biography == "" || newUser.City == "" || newUser.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	birthdate, err := time.Parse("2006-01-02", newUser.Birthdate)
	if err != nil {
		http.Error(w, "Invalid birthdate format. Use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	var userID string
	err = db.DB.QueryRow(
		"INSERT INTO users (first_name, last_name, birthdate, biography, city, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		newUser.FirstName, newUser.LastName, birthdate, newUser.Biography, newUser.City, hashedPassword,
	).Scan(&userID)
	if err != nil {
		http.Error(w, "Error inserting user into db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"user_id": userID}
	json.NewEncoder(w).Encode(response)
}
