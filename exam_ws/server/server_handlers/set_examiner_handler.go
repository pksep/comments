package server_handlers

import (
	"exam_ws/structure_managers"
	"log"
	"net"
	"strings"
)

func SetExaminerHandler(conn net.Conn, msg string, sm *structure_managers.ConnectionSessionStateManager) bool {
	// Принимает msg форматом examiner/<hash>
	msg_parts := strings.Split(msg, "/")
	hash := msg_parts[1]
	if !sm.NewConnectionSessionState(conn, msg) {
		log.Printf("Сессия %s уже существует смотрю есть ли подключения", hash)
	}

	if sm.SetExaminerSockAddr(hash, conn.RemoteAddr().String()) {
		log.Printf("Экзаменатор установлен в сессию %s", hash)
		return true
	} else {
		log.Printf("WARNING: Экзаменатор в сессии %s уже есть", hash)
		conn.Write([]byte("WARNING: Examiner exsists"))
		return false
	}
}
