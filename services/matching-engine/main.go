package main

import (
	"os"
	"os/signal"
	"runtime"

	"github.com/gerins/log"

	"matching-engine/cmd"
	"matching-engine/config"
	"matching-engine/internal/app"
)

func init() {
	log.Init()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		cfg  = config.ParseConfigFile("config.yaml")
		http = cmd.NewHttpServer(cfg)
		grpc = cmd.NewGRPCServer(cfg)
	)

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
