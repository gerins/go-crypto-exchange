package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gerins/log"
	middlewareLog "github.com/gerins/log/middleware/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"core-engine/config"
)

type HTTPServer struct {
	Server *echo.Echo
	cfg    *config.Config
}

// NewHTTPServer returns new HttpServer.
func NewHTTPServer(cfg *config.Config) *HTTPServer {
	return &HTTPServer{
		cfg:    cfg,
		Server: echo.New(),
	}
}

func (s *HTTPServer) Run() chan bool {
	// Apply middleware
	s.Server.Use(middlewareLog.Recover())
	s.Server.Use(middlewareLog.SetLogRequest()) // Mandatory
	s.Server.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: middlewareLog.SaveLogRequest(),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/"
		},
	}))

	// Init app
	s.Server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Start server
	go func() {
		s.Server.HideBanner = true
		address := fmt.Sprintf("%v:%v", s.cfg.App.HTTP.Host, s.cfg.App.HTTP.Port)
		if err := s.Server.Start(address); err != nil {
			// ErrServerClosed is expected behavior when exiting app
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("%v server, %v", s.cfg.App.Name, err)
			}
			log.Infof("%v server, %v", s.cfg.App.Name, err)
		}
	}()

	serverExitSignal := make(chan bool)
	go func() {
		<-serverExitSignal
		log.Info("stopping http server")
		if err := s.Server.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed stopping server, %v", err)
		}
		log.Info("finished stopping http server")
		serverExitSignal <- true // Send signal already finish the job
	}()

	return serverExitSignal
}
