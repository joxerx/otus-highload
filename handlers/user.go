package handlers

import (
	"encoding/json"
	"net/http"

	"otus-highload/db"
	"otus-highload/models"

	"github.com/gorilla/mux"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user models.User
	rows, err := db.ExecuteReadQuery("SELECT id, first_name, last_name, birthdate, biography, city FROM users WHERE id = $1", userID)
	if err != nil {
		http.Error(w, "Error executing read query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			http.Error(w, "Error reading user data", http.StatusInternalServerError)
			return
		}
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Birthdate, &user.Biography, &user.City)
	if err != nil {
		http.Error(w, "Error parsing user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
