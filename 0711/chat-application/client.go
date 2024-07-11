package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    // 連接到伺服器，假設伺服器位於 localhost 的 8080 端口
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
    defer conn.Close()

    // 開啟讀取伺服器訊息的 goroutine
    go readMessages(conn)

    // 從用戶端讀取輸入並發送給伺服器
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()
        _, err := fmt.Fprintf(conn, text+"\n")
        if err != nil {
            fmt.Println("Error sending message:", err)
            return
        }
    }
}

// 處理接收伺服器訊息
func readMessages(conn net.Conn) {
    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading message:", err)
            return
        }
        fmt.Print(msg)
    }
}