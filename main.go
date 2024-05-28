package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rest-api-gorilla/database"
	"rest-api-gorilla/handlers"
	"rest-api-gorilla/websocket"
)

func main() {
	database.InitDatabase()

	router := mux.NewRouter()

	router.HandleFunc("/api/messages", handlers.CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", handlers.GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{id}", handlers.UpdateMessage).Methods("PUT")
	router.HandleFunc("/api/messages/{id}", handlers.GetMessage2).Methods("GET")
	router.HandleFunc("/api/messages/{id}", handlers.DeleteMessage).Methods("DELETE")
	router.HandleFunc("/ws", websocket.HandleConnections)

	go websocket.HandleMessages()

	log.Println("Server started at :8000")
	server := http.ListenAndServe(":8000", router)
	log.Fatal(server)
}
