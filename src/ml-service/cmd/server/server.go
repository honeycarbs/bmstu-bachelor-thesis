package server

import (
	"errors"
	"fmt"
	"ml/internal/config"
	"ml/pkg/logging"
	"ml/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func Run(cfg *config.Config, handler http.Handler, logger logging.Logger) {
	var (
		s        Server
		listener net.Listener
	)

	s.httpServer = &http.Server{
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	logger.Infof("trying to listen to %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Listen.Port))
	if err != nil {
		logger.Fatal(err)
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		s.httpServer)

	logger.Println("application initialized and started")

	if err := s.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
