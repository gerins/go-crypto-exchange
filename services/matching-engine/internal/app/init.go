package app

import (
	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	"matching-engine/config"
	"matching-engine/internal/app/controller"
	"matching-engine/internal/app/usecase"
	"matching-engine/pkg/kafka"
	"matching-engine/pkg/redis"
)

func Init(e *echo.Echo, g *grpc.Server, cfg *config.Config) chan bool {
	var (
		exitSignal            = make(chan bool)
		validator             = validator.New()
		cache                 = redis.Init(cfg.Dependencies.Cache)
		kafkaConsumer         = kafka.NewConsumer(cfg.Dependencies.MessageBroker)
		kafkaProducer, writer = kafka.NewProducer(cfg.Dependencies.MessageBroker.Brokers)
	)

	// Init http router
	orderBookUsecase := usecase.NewOrderBook(cfg.Dependencies.MessageBroker.Consumer.Topic, cfg.Dependencies.MessageBroker.Producer.Topic, cache, kafkaProducer, validator)
	controller.NewHTTPHandler(orderBookUsecase, cfg.App.CtxTimeout).InitRoutes(e)
	controller.NewQueueHandler(kafkaConsumer, orderBookUsecase, cfg.App.CtxTimeout).StartConsumer()

	// Gracefull shutdown
	go func() {
		<-exitSignal // Receive exit signal
		log.Info("disconnecting service dependencies")

		if err := writer.Close(); err != nil {
			log.Error(err)
		}

		log.Info("finished disconnecting service dependencies")
		exitSignal <- true // Send signal already finish the job
	}()

	return exitSignal
}
