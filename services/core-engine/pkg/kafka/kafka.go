//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

//counterfeiter:generate . Producer
type Producer interface {
	Send(ctx context.Context, topic, key, message string) error
}

type producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers string) (Producer, *kafka.Writer) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      strings.Split(brokers, ","), //
		Balancer:     &kafka.Murmur2Balancer{},    // Partition balancer
		MaxAttempts:  3,                           // Limit on how many attempts will be made to deliver a message.
		BatchTimeout: time.Second,                 // Time limit on how often incomplete message batches will be flushed to kafka.
		RequiredAcks: int(kafka.RequireOne),       // Wait for all replicas
	})

	return &producer{writer: writer}, writer
}

// Send is used for sending message to Kafka
func (kp *producer) Send(ctx context.Context, topic, key, message string) error {
	newMessage := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: []byte(message),
	}

	// Sending message
	if err := kp.writer.WriteMessages(ctx, newMessage); err != nil {
		return err
	}

	return nil
}
