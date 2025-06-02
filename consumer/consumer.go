package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
)

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

	chConn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"clickhouse-server:9000"},
		Auth: clickhouse.Auth{
			Database: "marketplace_analytics",
			Username: "default",
			Password: "",
		},
		Debug: false,
	})
	if err != nil {
		log.Fatalf("clickhouse connection error: %v", err)
	}
	defer chConn.Close()

	ctx := context.Background()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
		MaxWait:     1 * time.Second,
	})
	defer reader.Close()

	fmt.Println("Kafka consumer started. Listening for events...")

	for {
		msg, err := reader.ReadMessage(ctx)
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

		// Парсим timestamp в time.Time
		eventTime, err := time.Parse(time.RFC3339, event.Timestamp)
		if err != nil {
			log.Printf("error while parsing timestamp: %v\n", err)
			eventTime = time.Now()
		}

		// Пишем событие в ClickHouse
		batch, err := chConn.PrepareBatch(ctx, "INSERT INTO marketplace_analytics.user_events (event_time, url, action, ids)")
		if err != nil {
			log.Printf("clickhouse prepare batch error: %v", err)
			continue
		}

		err = batch.Append(eventTime, event.URL, event.Action, event.Ids)
		if err != nil {
			log.Printf("clickhouse append batch error: %v", err)
			continue
		}

		err = batch.Send()
		if err != nil {
			log.Printf("clickhouse send batch error: %v", err)
			continue
		}

		log.Printf("✅ received event: URL=%s, action=%s, timestamp=%s\n", event.URL, event.Action, event.Timestamp)
	}
}
