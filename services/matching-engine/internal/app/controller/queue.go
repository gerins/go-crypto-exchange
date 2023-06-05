package controller

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gerins/log"
	"github.com/segmentio/kafka-go"

	"matching-engine/internal/app/model"
)

type queueHandler struct {
	kafkaConsumer *kafka.Reader
	engine        model.Engine
	timeout       time.Duration
}

func NewQueueHandler(kafkaConsumer *kafka.Reader, processor model.Engine, timeout time.Duration) *queueHandler {
	return &queueHandler{
		kafkaConsumer: kafkaConsumer,
		engine:        processor,
		timeout:       timeout,
	}
}

func (h *queueHandler) StartConsumer() {
	go func() {
		for {
			kafkaMessage, err := h.kafkaConsumer.FetchMessage(context.Background())
			if err != nil {
				continue
			}

			ctx := context.TODO()
			if err := h.OrderHandler(ctx, kafkaMessage.Value); err != nil {
				log.Context(ctx).Error(err)
				return
			}

			// Commit message
			h.kafkaConsumer.CommitMessages(ctx, kafkaMessage)
		}
	}()
}

func (h *queueHandler) OrderHandler(ctx context.Context, msg []byte) error {
	var payload model.Order
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Context(ctx).Error(err)
		return err
	}

	if err := h.engine.Execute(ctx, payload); err != nil {
		return err
	}

	return nil
}
