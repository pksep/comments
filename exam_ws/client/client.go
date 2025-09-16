package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	log.Print("Попытка подключения к серверу...")

	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Printf("Ошибка при подключение к серверу localhost:8080: %s", err)
	}

	defer conn.Close()

	// Фоновое чтение
	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n')

			if err != nil {
				log.Fatalf("Соединение потеряно: %s", err)
			}

			log.Printf("Сервер: %s", msg)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()

		if text == "ex" {
			log.Print("Выход...")
			break
		}

		if text == "" {
			continue
		}

		byte_len, err := conn.Write([]byte(text + "\n"))

		if err != nil {
			log.Fatalf("Критичная ошибка при отправке сообщения: %s", err)
		}

		log.Printf("Успешная отправка сообщения длинной байтов: %d", byte_len)

		time.Sleep(100 * time.Millisecond)
	}

	err_scanner := scanner.Err()
	if err != nil {
		log.Printf("Ошибка чтения ввода: %s", err_scanner)
	}

}
