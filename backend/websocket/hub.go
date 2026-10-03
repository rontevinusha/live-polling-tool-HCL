package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)
var CurrentHub *Hub
type Hub struct {
	Clients map[*websocket.Conn]bool
	Mutex   sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	h.Clients[conn] = true
}

func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	delete(h.Clients, conn)
	conn.Close()
}

func (h *Hub) Broadcast(message interface{}) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	for conn := range h.Clients {
		err := conn.WriteJSON(message)
		if err != nil {
			delete(h.Clients, conn)
			conn.Close()
		}
	}
}