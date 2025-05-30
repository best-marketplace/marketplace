package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		// topic: topic,
		writer: &kafka.Writer{
			Addr: kafka.TCP(brokers...),
			// Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (kp *KafkaProducer) Send(ctx context.Context, event any, topic string) error {
	msgBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	fmt.Println(topic)
	message := kafka.Message{
		Topic: topic,
		Key:   nil,
		Value: msgBytes,
		Time:  time.Now(),
	}

	return kp.writer.WriteMessages(ctx, message)
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}
