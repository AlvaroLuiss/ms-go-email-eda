package main

import (
	"fmt"
	"log"
	consumer "ms-go-eda-email/kafka"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Nenhum arquivo .env encontrado")
	}

	fmt.Println("[app] Iniciando consumidor de eventos de usuário...")
	consumer.StartConsumer()
}
