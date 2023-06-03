package app

import (
	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	"core-engine/config"
	"core-engine/internal/app/domains/user"
	"core-engine/pkg/gorm"
	"core-engine/pkg/kafka"
	"core-engine/pkg/redis"
)

func Init(e *echo.Echo, g *grpc.Server, cfg *config.Config) chan bool {
	var (
		exitSignal     = make(chan bool)
		validator      = validator.New()
		defaultTimeout = cfg.App.CtxTimeout
		readDatabase   = gorm.Init(cfg.Dependencies.Database.Read)
		writeDatabase  = gorm.Init(cfg.Dependencies.Database.Write)
		_              = redis.Init(cfg.Dependencies.Cache)
		_, writer      = kafka.NewProducer(cfg.Dependencies.MessageBroker.Brokers)
	)

	// Init http router
	userRepository := user.NewRepository(readDatabase, writeDatabase)
	userUsecase := user.NewUsecase(cfg.Security, validator, userRepository)
	user.NewHTTPHandler(userUsecase, defaultTimeout).InitRoutes(e)

	// Graceful shutdown
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
