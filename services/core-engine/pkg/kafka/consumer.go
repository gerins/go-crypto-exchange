package kafka

import (
	"strings"
	"time"

	"github.com/segmentio/kafka-go"

	"core-engine/config"
)

func NewConsumer(cfg config.MessageBroker, topic string) *kafka.Reader {
	consumerConfig := kafka.ReaderConfig{
		Brokers:         strings.Split(cfg.Brokers, ","), // "localhost:9092,localhost:9092"
		GroupID:         cfg.Group,
		Topic:           topic,
		MinBytes:        10e3, // 10KB
		MaxBytes:        10e6, // 10MB
		MaxWait:         10 * time.Millisecond,
		ReadLagInterval: -1,
		CommitInterval:  time.Second,
		StartOffset:     kafka.LastOffset,
	}

	reader := kafka.NewReader(consumerConfig)
	return reader
}
