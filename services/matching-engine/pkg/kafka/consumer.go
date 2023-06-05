package kafka

import (
	"strings"
	"time"

	"github.com/segmentio/kafka-go"

	"matching-engine/config"
)

func NewConsumer(cfg config.MessageBroker) *kafka.Reader {
	consumerConfig := kafka.ReaderConfig{
		Brokers:         strings.Split(cfg.Brokers, ","), // "localhost:9092,localhost:9092"
		GroupID:         cfg.Group,
		Topic:           cfg.Consumer.Topic,
		MinBytes:        10e3, // 10KB
		MaxBytes:        10e6, // 10MB
		MaxWait:         50 * time.Millisecond,
		ReadLagInterval: -1,
		CommitInterval:  50 * time.Millisecond,
		StartOffset:     kafka.LastOffset,
	}

	reader := kafka.NewReader(consumerConfig)
	return reader
}
