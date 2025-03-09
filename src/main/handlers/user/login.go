package user

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
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rows, err := db.ExecuteReadQuery("SELECT password FROM users WHERE id = $1", loginReq.ID)
	if err != nil {
		http.Error(w, "Invalid user ID or password", http.StatusUnauthorized)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		http.Error(w, "Invalid user ID or password", http.StatusUnauthorized)
		return
	}

	var hashedPassword string
	if err := rows.Scan(&hashedPassword); err != nil {
		http.Error(w, "Error scanning password", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password)); err != nil {
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
