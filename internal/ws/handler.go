package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) Connect(c *gin.Context) {
	role := c.Query("role")
	roomID := c.Query("room")

	if role == "" || roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role and room required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		Conn: conn,
		Role: role,
		Send: make(chan []byte, 256),
	}

	room := h.hub.GetOrCreateRoom(roomID)
	if !room.AddClient(client) {
		conn.WriteMessage(websocket.TextMessage, []byte("❌ Room is full"))
		conn.Close()
		return
	}

	// стартуем горутины
	go clientWriter(client)
	go clientReader(client, room)
}

func clientWriter(c *Client) {
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func clientReader(c *Client, room *Room) {
	defer c.Conn.Close()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		if c.Role == "examer" {
			log.Println("Examer ответил:", string(msg))
			// когда экзаменатор ответил — шлём следующий вопрос
			room.NextQuestion()
		}
	}
}
