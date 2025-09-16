package server_handlers

import (
	"exam_ws/structure_managers"
	"log"
	"net"
)

func SetUserHandler(conn net.Conn, msg string, sm *structure_managers.ConnectionSessionStateManager) bool {
	// Принимает msg форматом <hash>
	if !sm.NewConnectionSessionState(conn, msg) {
		log.Printf("Сессия %s уже существует смотрю есть ли подключения", msg)
	}

	if sm.SetUserSockAddr(msg, conn.RemoteAddr().String()) {
		log.Printf("Пользователь установлен в сессию %s", msg)
		return true
	} else {
		log.Printf("WARNING: Пользователь в сессии %s уже есть", msg)
		conn.Write([]byte("WARNING: User exsists"))
		return false
	}
}
