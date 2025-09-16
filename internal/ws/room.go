package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID      string
	Examer  *Client
	Student *Client
	timer   *time.Timer
	mu      sync.Mutex
}

func (r *Room) AddClient(c *Client) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch c.Role {
	case "examer":
		if r.Examer != nil {
			return false
		}
		r.Examer = c
	case "student":
		if r.Student != nil {
			return false
		}
		r.Student = c
	}

	// если оба подключились — стартуем таймер
	if r.Examer != nil && r.Student != nil && r.timer == nil {
		r.timer = time.AfterFunc(time.Minute, func() {
			r.Broadcast([]byte("⏳ Время истекло"))
		})
		r.Broadcast([]byte("✅ Сессия началась"))
	}
	return true
}

func (r *Room) Broadcast(msg []byte) {
	if r.Examer != nil {
		r.Examer.Send <- msg
	}
	if r.Student != nil {
		r.Student.Send <- msg
	}
}

func (r *Room) NextQuestion() {
	if r.Student != nil {
		r.Student.Send <- []byte(fmt.Sprintf("Следующий вопрос: %s", uuid.NewString()))
	}
}
