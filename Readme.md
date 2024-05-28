# CRUD Golang (REST API, WEBSOCKET, SQLITE)

## Struktur Proyek

```go
myapp/
├── main.go
├── handlers/
│   └── handlers.go
├── models/
│   └── message.go
├── database/
│   └── db.go
├── websocket/
│   └── websocket.go

```


## Inisialisasi Proyek

```bash
mkdir rest-api-gorilla
cd rest-api-gorilla

go mod init rest-api-gorilla


go get github.com/gorilla/mux
go get github.com/gorilla/websocket
go get modernc.org/sqlite


```


## Konfigurasi Database

Buat file d`atabase/db.go` untuk mengatur koneksi ke **SQLite**:

```go
package database

import (
    "database/sql"
    "log"

    _ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDatabase() {
    var err error
    DB, err = sql.Open("sqlite", "./messages.db")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    createTable()
}

func createTable() {
    query := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL
    );`
    _, err := DB.Exec(query)
    if err != nil {
        log.Fatalf("Failed to create table: %v", err)
    }
}

```


## Mendefinisikan Model

Buat file `models/message.go` untuk model **Message**:

```go
package models

type Message struct {
    ID      int    `json:"id"`
    Content string `json:"content"`
}

```

## Implementasi Handlers

Buat file `handlers/handlers.go` untuk mengatur routing dan logika API:

```go
package handlers

import (
    "encoding/json"
    "log"
    "myapp/database"
    "myapp/models"
    "myapp/websocket"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
    var msg models.Message
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    query := "INSERT INTO messages (content) VALUES (?)"
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
    rows, err := database.DB.Query("SELECT id, content FROM messages")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var messages []models.Message
    for rows.Next() {
        var msg models.Message
        if err := rows.Scan(&msg.ID, &msg.Content); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        messages = append(messages, msg)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
    var msg models.Message
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    query := "UPDATE messages SET content = ? WHERE id = ?"
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

    query := "DELETE FROM messages WHERE id = ?"
    _, err := database.DB.Exec(query, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

```

## Konfigurasi WebSocket

Buat file `websocket/websocket.go` untuk mengatur **WebSocket**:

```go
package websocket

import (
    "encoding/json"
    "log"
    "myapp/models"
    "net/http"

    "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)
var upgrader = websocket.Upgrader{}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatalf("Failed to upgrade connection: %v", err)
    }
    defer ws.Close()

    clients[ws] = true

    for {
        var msg models.Message
        err := ws.ReadJSON(&msg)
        if err != nil {
            delete(clients, ws)
            break
        }
    }
}

func HandleMessages() {
    for {
        msg := <-broadcast
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("Error sending message: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}

func BroadcastMessage(msg models.Message) {
    broadcast <- msg
}

```


## Inisialisasi Server

Buat file `main.go` untuk menginisialisasi server dan mengatur **routing**:

```go
package main

import (
    "log"
    "myapp/database"
    "myapp/handlers"
    "myapp/websocket"
    "net/http"

    "github.com/gorilla/mux"
)

func main() {
    database.InitDatabase()

    router := mux.NewRouter()

    router.HandleFunc("/api/messages", handlers.CreateMessage).Methods("POST")
    router.HandleFunc("/api/messages", handlers.GetMessages).Methods("GET")
    router.HandleFunc("/api/messages/{id}", handlers.UpdateMessage).Methods("PUT")
    router.HandleFunc("/api/messages/{id}", handlers.DeleteMessage).Methods("DELETE")
    router.HandleFunc("/ws", websocket.HandleConnections)

    go websocket.HandleMessages()

    log.Println("Server started at :8000")
    log.Fatal(http.ListenAndServe(":8000", router))
}

```

## Menjalankan Aplikasi

1. Pastikan semua file sudah dibuat dengan benar.
2. Jalankan aplikasi dengan perintah berikut:

```bash
go run main.go
```