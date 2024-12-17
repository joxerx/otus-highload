package user

import (
	"encoding/json"
	"log"
	"net/http"
	"otus-highload/db"
	"otus-highload/models"
)

func SearchUserHandler(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	if firstName == "" {
		http.Error(w, "first_name is required", http.StatusBadRequest)
		return
	}

	if lastName == "" {
		http.Error(w, "last_name is required", http.StatusBadRequest)
		return
	}

	query := "SELECT id, first_name, last_name, birthdate, biography, city FROM users WHERE first_name LIKE $1 AND last_name LIKE $2 ORDER BY id"
	rows, err := db.ExecuteReadQuery(query, firstName+"%", lastName+"%")
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Birthdate, &user.Biography, &user.City); err != nil {
			log.Println("Error scanning user:", err)
			http.Error(w, "Error parsing users from db", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error processing rows:", err)
		http.Error(w, "Error processing users from db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error sending response", http.StatusInternalServerError)
		return
	}
}
