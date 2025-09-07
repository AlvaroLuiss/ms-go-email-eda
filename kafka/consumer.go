package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"

	"ms-go-eda-email/email"
)

type UserRegisteredEvent struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type UserLoginEvent struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	At     string `json:"at"`
}

func StartConsumer() {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	groupID := getEnv("KAFKA_GROUP_ID", "go-event-consumer-group")

	// Consumer para user.registered
	// Usando goroutine para rodar em paralelo
	go consumeUserRegistered(broker, groupID)

	// Consumer para user.login
	consumeUserLogin(broker, groupID)
}

func consumeUserRegistered(broker, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       "user.registered",
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
	})
	defer r.Close()

	fmt.Println("[kafka] Consumindo tópico: user.registered")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Erro lendo mensagem user.registered:", err)
			continue
		}

		var event UserRegisteredEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			fmt.Println("Erro parseando JSON user.registered:", err)
			continue
		}

		fmt.Printf("[kafka] Usuário registrado: %+v\n", event)

		// Aqui envia o email de boas vindas
		if err := email.SendWelcomeEmail(event.Email, event.ID); err != nil {
			fmt.Println("Erro enviando email de boas-vindas:", err)
		}


		fmt.Printf("[audit] Usuário %s (%s) registrado em %s\n", event.ID, event.Email, event.CreatedAt)
	}
}

func consumeUserLogin(broker, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       "user.login",
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
	})
	defer r.Close()

	fmt.Println("[kafka] Consumindo tópico: user.login")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Erro lendo mensagem user.login:", err)
			continue
		}

		var event UserLoginEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			fmt.Println("Erro parseando JSON user.login:", err)
			continue
		}

		fmt.Printf("[kafka] Login realizado: %+v\n", event)

		fmt.Printf("[audit] Login do usuário %s (%s) em %s\n", event.UserID, event.Email, event.At)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}