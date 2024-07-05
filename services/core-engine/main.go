package main

import (
	"os"
	"os/signal"
	"runtime"

	"github.com/gerins/log"

	"core-engine/cmd"
	"core-engine/config"
	"core-engine/internal/app"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		cfg  = config.ParseConfigFile("config.yaml")
		http = cmd.NewHTTPServer(cfg)
		grpc = cmd.NewGRPCServer(cfg)
	)

	// Init log
	log.InitWithConfig(log.Config{
		LogToTerminal:     cfg.App.Logging.LogToTerminal,
		LogToFile:         cfg.App.Logging.LogToFile,
		Location:          cfg.App.Logging.Location,
		FileLogName:       cfg.App.Logging.FileLogName,
		MaxAge:            cfg.App.Logging.MaxAge,
		RotationFile:      cfg.App.Logging.RotationFile,
		HideSensitiveData: cfg.App.Logging.HideSensitiveData,
	})

	// Init app
	appExitSignal := app.Init(http.Server, grpc.Server, cfg)

	// Run server
	httpExitSignal := http.Run()
	grpcExitSignal := grpc.Run()

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)
	for range interruptSignal {
		appExitSignal <- true
		httpExitSignal <- true
		grpcExitSignal <- true
		<-appExitSignal  // Finish
		<-httpExitSignal // Finish
		<-grpcExitSignal // Finish
		return           // Now we can safely exit the app
	}
}
