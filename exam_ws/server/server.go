package main

import (
	"exam_ws/server_handlers"
	"exam_ws/structure_managers"
	"log"
	"net"
)

func main() {
	log.Print("Запуск сервера на порте 8080 TCP")

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}

	defer listener.Close() // После завершения функции закрыть прослушку tcp:8080

	log.Println("Сервер запущен на localhost:8080")
	log.Println("Ожидание подключений...")

	conn_session_state_manager := structure_managers.NewConnectionSessionStateManager()

	// Бесконечный цикл с ожиданием сообщений
	for {
		conn, err := listener.Accept() // Ожидаем нового подключени БЛОКИРУЕТ ВЫПОЛНЕНИЕ
		client_addr := conn.RemoteAddr().String()

		if err != nil {
			log.Printf("Ошибка подключения с клиентом %s", client_addr)
			continue
		}

		go server_handlers.MainHandlerConnection(conn, conn_session_state_manager)

	}

}
