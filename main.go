package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // для разработки
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("New HTTP request:", r.RemoteAddr)

	// обновление соединения до WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	defer ws.Close()

	log.Println("WebSocket connected:", r.RemoteAddr)

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", message)

		prefix := []byte("Echo: ")
		response := append(prefix, message...)

		if err := ws.WriteMessage(messageType, response); err != nil {
			log.Println("Write error:", err)
			break
		}
	}

	log.Println("Client disconnected:", r.RemoteAddr)
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
