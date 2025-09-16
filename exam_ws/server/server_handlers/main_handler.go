package server_handlers

import (
	"bufio"
	"exam_ws/structure_managers"
	"log"
	"net"
	"strings"
)

func MainHandlerConnection(conn net.Conn, sm *structure_managers.ConnectionSessionStateManager) {
	reader := bufio.NewReader(conn)

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Новый клиент подключился: %s\n", clientAddr)

	var user_role string
	var hash string

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Клиент %s отключился\n", clientAddr)
			return
		}

		msg = strings.TrimSpace(msg)

		if user_role == "" {
			if strings.Contains(msg, "/") {
				msg_parts := strings.Split(msg, "/")
				hash = msg_parts[1]
				user_role = "examiner"
			} else {
				hash = msg
				user_role = "user"
			}
			defer sm.RemoveConnect(conn, hash, user_role)
		}

		if strings.Contains(msg, "examiner") {
			SetExaminerHandler(conn, msg, sm)
		} else {
			SetUserHandler(conn, msg, sm)
		}
	}
}
