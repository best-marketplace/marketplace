package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		topic: topic,
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (kp *KafkaProducer) Send(ctx context.Context, event any) error {
	msgBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   nil,
		Value: msgBytes,
		Time:  time.Now(),
	}

	return kp.writer.WriteMessages(ctx, message)
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}
