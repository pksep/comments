package ws

import "sync"

type Hub struct {
	rooms map[string]*Room
	mu    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

func (h *Hub) GetOrCreateRoom(id string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.rooms[id]
	if !ok {
		room = &Room{ID: id}
		h.rooms[id] = room
	}
	return room
}
