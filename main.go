package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"rest-api-gorilla/database"
	"rest-api-gorilla/handlers"
	"rest-api-gorilla/websocket"
)

var pathEnv = ".env"

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
	var port = "8000"
	err := handlers.LoadEnv(pathEnv)
	if err != nil {
		log.Println("Load env :", err)
	}

	port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("Server started at :" + port)
	server := http.ListenAndServe(":"+port, router)
	log.Fatal(server)
}
