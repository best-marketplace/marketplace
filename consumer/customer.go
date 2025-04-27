package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Event соответствует структуре, которую ты отправляешь
type Event struct {
	URL       string   `json:"url"`
	Ids       []string `json:"ids"`
	Action    string   `json:"action"`
	Timestamp string   `json:"timestamp"`
}

func main() {
	brokers := []string{"kafka:9092"}
	topic := "user-events"
	groupID := "event-logger-group"

	// Создаём нового reader'а
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,    // минимальный размер сообщения
		MaxBytes:    10e6, // максимум — 10MB
		MaxWait:     1 * time.Second,
	})
	defer reader.Close()

	fmt.Println("Kafka consumer started. Listening for events...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error while reading message: %v\n", err)
			continue
		}

		var event Event
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("error while decoding message: %v\n", err)
			continue
		}

		log.Printf("✅ received event: URL=%s, action=%s\n,timestamp:=%s\n", event.URL, event.Action, event.Timestamp)
		log.Println(event.Ids)
	}
}
