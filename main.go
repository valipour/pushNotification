package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	username  = "کاربر"
	mu        sync.RWMutex
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func broadcast(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("خطا در ارسال:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

func setNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "فقط POST مجاز است", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "پارامتر 'name' الزامی است", http.StatusBadRequest)
		return
	}

	mu.Lock()
	username = name
	mu.Unlock()

	msg := fmt.Sprintf("نام جدید تنظیم شد: %s", name)
	go broadcast(msg)

	fmt.Fprintln(w, msg)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("خطا در WebSocket:", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

		mu.RLock()
		name := username
		mu.RUnlock()

		msg := fmt.Sprintf("سلام %s! نوتیف #%d - %s", name, i+1, time.Now().Format("15:04:05"))
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println("خطا در ارسال نوتیف:", err)
			break
		}
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/set-name", setNameHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("سرور روی http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
