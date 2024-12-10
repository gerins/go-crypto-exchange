package app

import (
	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	"core-engine/config"
	"core-engine/internal/app/domains/handler"
	"core-engine/internal/app/domains/repository"
	"core-engine/internal/app/domains/usecase"
	"core-engine/pkg/gorm"
	"core-engine/pkg/kafka"
	"core-engine/pkg/redis"
)

func Init(e *echo.Echo, g *grpc.Server, cfg *config.Config) chan bool {
	var (
		exitSignal         = make(chan bool)
		validator          = validator.New()
		apiTimeout         = cfg.App.HTTP.CtxTimeout
		redisCache         = redis.Init(cfg.Dependencies.Cache)
		redisLock          = redis.InitLock(redisCache)
		readDatabase       = gorm.InitPostgres(cfg.Dependencies.Database.Read)
		writeDatabase      = gorm.InitPostgres(cfg.Dependencies.Database.Write)
		matchOrderConsumer = kafka.NewConsumer(cfg.Dependencies.MessageBroker, cfg.Dependencies.MessageBroker.Consumer.Topic.MatchOrder)
		producer, writer   = kafka.NewProducer(cfg.Dependencies.MessageBroker.Brokers)
	)

	// Repository
	userRepository := repository.NewUserRepository(readDatabase, writeDatabase)
	orderRepository := repository.NewOrderRepository(readDatabase, writeDatabase)

	// Usecase
	userUsecase := usecase.NewUserUsecase(cfg.Security, validator, userRepository)
	orderUsecase := usecase.NewOrderUsecase(redisLock, writeDatabase, producer, validator, orderRepository, userRepository)

	// Handler
	handler.NewUserHandler(userUsecase, apiTimeout).InitRoutes(e)
	handler.NewOrderHTTPHandler(orderUsecase, apiTimeout, cfg.Security).InitRoutes(e)
	handler.NewOrderQueueHandler(matchOrderConsumer, orderUsecase, apiTimeout).StartConsumer()

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

		if err := redisCache.Close(); err != nil {
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
