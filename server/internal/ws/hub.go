package ws

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"hanunted-house/server/internal/room"
)

type Hub struct {
	manager *room.Manager
	logger  *slog.Logger
	mu      sync.RWMutex
	rooms   map[string]map[string]*websocket.Conn
}

func NewHub(manager *room.Manager, logger *slog.Logger) *Hub {
	return &Hub{manager: manager, logger: logger, rooms: map[string]map[string]*websocket.Conn{}}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (h *Hub) Handle(w http.ResponseWriter, r *http.Request, roomCode, playerID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("upgrade failed", "error", err)
		return
	}

	h.addConnection(roomCode, playerID, conn)
	h.broadcastRoom(roomCode)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.removeConnection(roomCode, playerID)
			h.manager.Disconnect(roomCode, playerID)
			h.broadcastRoom(roomCode)
			_ = conn.Close()
			return
		}
	}
}

func (h *Hub) broadcastRoom(roomCode string) {
	r, err := h.manager.GetRoomByCode(roomCode)
	if err != nil {
		return
	}

	h.mu.RLock()
	clients := h.rooms[roomCode]
	h.mu.RUnlock()
	for playerID, conn := range clients {
		view := r.PlayerView(playerID)
		if err := conn.WriteJSON(map[string]any{"type": "snapshot", "payload": view}); err != nil {
			h.logger.Warn("broadcast failed", "error", err)
		}
	}
}

func (h *Hub) Broadcast(roomCode string) { h.broadcastRoom(roomCode) }

func (h *Hub) addConnection(roomCode, playerID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[roomCode]; !ok {
		h.rooms[roomCode] = map[string]*websocket.Conn{}
	}
	h.rooms[roomCode][playerID] = conn
}

func (h *Hub) removeConnection(roomCode, playerID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[roomCode]; !ok {
		return
	}
	delete(h.rooms[roomCode], playerID)
}
