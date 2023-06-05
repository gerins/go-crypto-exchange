//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
package kafka

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gerins/log"
	"github.com/segmentio/kafka-go"
)

//counterfeiter:generate . Producer
type Producer interface {
	Send(ctx context.Context, topic, key string, payload interface{}) error
}

type producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers string) (Producer, *kafka.Writer) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      strings.Split(brokers, ","), //
		Balancer:     &kafka.Murmur2Balancer{},    // Partition balancer
		MaxAttempts:  3,                           // Limit on how many attempts will be made to deliver a message.
		BatchTimeout: 1 * time.Millisecond,        // Time limit on how often incomplete message batches will be flushed to kafka.
		RequiredAcks: int(kafka.RequireOne),       // Wait for all replicas
	})

	return &producer{writer: writer}, writer
}

// Send is used for sending message to Kafka
func (kp *producer) Send(ctx context.Context, topic, key string, payload interface{}) error {
	defer log.Context(ctx).RecordDuration("kafka publisher").Stop()

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Context(ctx).Error(err)
		return err
	}

	newMessage := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payloadJSON,
	}

	// Sending message
	if err := kp.writer.WriteMessages(ctx, newMessage); err != nil {
		return err
	}

	return nil
}
