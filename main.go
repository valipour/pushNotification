package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

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

// ارسال پیام به همه کلاینت‌های متصل
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

// تغییر نام و ارسال نوتیف
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
	go broadcast(msg) // فقط همینجا نوتیف می‌فرستیم

	fmt.Fprintln(w, msg)
}

// مدیریت اتصال WebSocket
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

	// اتصال رو باز نگه می‌داریم ولی منتظر پیام از سمت کلاینت نیستیم
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
}

func main() {
	http.HandleFunc("/set-name", setNameHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("سرور روی http://localhost:8080 اجرا شد")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
