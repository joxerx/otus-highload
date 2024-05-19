package handlers

import (
	"encoding/json"
	"net/http"

	"otus-highload/db"
	"otus-highload/models"
	"otus-highload/utils"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err = db.DB.QueryRow("SELECT password FROM users WHERE id = $1", loginReq.ID).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "Invalid user ID or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password))
	if err != nil {
		http.Error(w, "Invalid user ID or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(loginReq.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"auth_token": token}
	json.NewEncoder(w).Encode(response)
}
