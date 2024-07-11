package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
    Username string `json:"username"`
    Message  string `json:"message"`
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Error while upgrading connection:", err)
        return
    }
    fmt.Println("Client connected")

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Error while reading message:", err)
            break
        }

        var message Message
        err = json.Unmarshal(msg, &message)
        if err != nil {
            fmt.Println("Error while unmarshaling message:", err)
            continue
        }

        fmt.Printf("Received message from %s: %s\n", message.Username, message.Message)

        // 將消息回傳給 Django
        err = conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            fmt.Println("Error while writing message:", err)
            break
        }
    }
}

func main() {
    http.HandleFunc("/ws/chat/", handleConnection)
    fmt.Println("Chat API server started at ws://localhost:8080/ws/chat/")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}
