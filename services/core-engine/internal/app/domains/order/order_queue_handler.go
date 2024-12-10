package order

import (
	"context"
	"fmt"
	"time"

	"github.com/gerins/log"
	"github.com/segmentio/kafka-go"

	"core-engine/internal/app/domains/order/model"
)

type queueHandler struct {
	kafkaConsumer *kafka.Reader
	orderUsecase  model.Usecase
	timeout       time.Duration
}

func NewQueueHandler(kafkaConsumer *kafka.Reader, orderUsecase model.Usecase, timeout time.Duration) *queueHandler {
	return &queueHandler{
		kafkaConsumer: kafkaConsumer,
		orderUsecase:  orderUsecase,
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

			go func() {
				logging := log.NewRequest()
				logging.Method = kafkaMessage.Topic
				logging.IP = string(kafkaMessage.Key)
				logging.URL = fmt.Sprintf("partition %v offset %v", kafkaMessage.Partition, kafkaMessage.Offset)

				// Parent context
				ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
				defer func() { logging.Save(); cancel() }()

				// Proceed match order
				if err := h.MatchOrderHandler(logging.SaveToContext(ctx), kafkaMessage.Value); err != nil {
					return // Dont commit message if error occur
				}

				// Commit message
				if err := h.kafkaConsumer.CommitMessages(ctx, kafkaMessage); err != nil {
					log.Context(ctx).Error(err)
					return
				}
			}()
		}
	}()
}

func (h *queueHandler) MatchOrderHandler(ctx context.Context, msg []byte) error {
	var payload model.BulkTradeRequest

	if err := payload.FromJSON(msg); err != nil {
		log.Context(ctx).Error(err)
		return err
	}

	log.Context(ctx).ReqBody = payload

	for i := 0; i < len(payload); i++ {
		if err := h.orderUsecase.MatchOrder(ctx, payload[i]); err != nil {
			return err
		}
	}

	return nil
}
