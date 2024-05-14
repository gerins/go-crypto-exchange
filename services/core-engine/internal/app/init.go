package app

import (
	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	"core-engine/config"
	"core-engine/internal/app/domains/order"
	"core-engine/internal/app/domains/user"
	"core-engine/pkg/gorm"
	"core-engine/pkg/kafka"
	"core-engine/pkg/redis"
)

func Init(e *echo.Echo, g *grpc.Server, cfg *config.Config) chan bool {
	var (
		exitSignal         = make(chan bool)
		validator          = validator.New()
		defaultTimeout     = cfg.App.HTTP.CtxTimeout
		readDatabase       = gorm.Init(cfg.Dependencies.Database.Read)
		writeDatabase      = gorm.Init(cfg.Dependencies.Database.Write)
		redis              = redis.Init(cfg.Dependencies.Cache)
		matchOrderConsumer = kafka.NewConsumer(cfg.Dependencies.MessageBroker, cfg.Dependencies.MessageBroker.Consumer.Topic.MatchOrder)
		producer, writer   = kafka.NewProducer(cfg.Dependencies.MessageBroker.Brokers)
	)

	// Init http router
	{
		// User Domain
		userRepository := user.NewRepository(readDatabase, writeDatabase)
		userUsecase := user.NewUsecase(cfg.Security, validator, userRepository)
		user.NewHTTPHandler(userUsecase, defaultTimeout).InitRoutes(e)

		// Order Domain
		orderRepository := order.NewRepository(readDatabase, writeDatabase)
		orderUsecase := order.NewUsecase(writeDatabase, producer, validator, orderRepository, userRepository)
		order.NewHTTPHandler(orderUsecase, defaultTimeout).InitRoutes(e, cfg.Security)
		order.NewQueueHandler(matchOrderConsumer, orderUsecase, defaultTimeout).StartConsumer()
	}

	// Graceful shutdown
	go func() {
		<-exitSignal // Receive exit signal
		log.Info("disconnecting service dependencies")

		if err := matchOrderConsumer.Close(); err != nil {
			log.Error(err)
		}

		if err := writer.Close(); err != nil {
			log.Error(err)
		}

		if err := redis.Close(); err != nil {
			log.Error(err)
		}

		if readDatabase, err := readDatabase.DB(); err == nil {
			if err = readDatabase.Close(); err != nil {
				log.Error(err)
			}
		}

		if writeDatabase, err := writeDatabase.DB(); err == nil {
			if err = writeDatabase.Close(); err != nil {
				log.Error(err)
			}
		}

		log.Info("finished disconnecting service dependencies")
		exitSignal <- true // Send signal already finish the job
	}()

	return exitSignal
}
