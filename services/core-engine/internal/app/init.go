package app

import (
	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	"core-engine/config"
	"core-engine/internal/app/controller"
	"core-engine/internal/app/usecase"
	"core-engine/pkg/kafka"
	"core-engine/pkg/redis"
)

func Init(e *echo.Echo, g *grpc.Server, cfg *config.Config) chan bool {
	var (
		exitSignal            = make(chan bool)
		validator             = validator.New()
		cache                 = redis.Init(cfg.Dependencies.Cache)
		kafkaProducer, writer = kafka.NewProducer(cfg.Dependencies.MessageBroker.Brokers)
	)

	// Init http router
	orderBookUsecase := usecase.NewOrderBook(validator, kafkaProducer, cache)
	controller.NewHTTPHandler(orderBookUsecase, cfg.App.CtxTimeout).InitRoutes(e)

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
