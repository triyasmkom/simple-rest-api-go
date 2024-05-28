package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"rest-api-gorilla/database"
	"rest-api-gorilla/models"
	"rest-api-gorilla/websocket"
	"strconv"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Mesaage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO messages (content) VALUES (?)`
	result, err := database.DB.Exec(query, msg.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	msg.ID = int(id)

	websocket.BroadcastMessage(msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`SELECT id, content FROM messages WHERE 1`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var message []models.Mesaage
	for rows.Next() {
		var msg models.Mesaage
		if err := rows.Scan(&msg.ID, &msg.Content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		message = append(message, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Mesaage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	query := `UPDATE messages SET content = ? WHERE id = ?`
	_, err := database.DB.Exec(query, msg.Content, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg.ID = id

	websocket.BroadcastMessage(msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)

}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	query := `DELETE FROM messages WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
