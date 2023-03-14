package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ThreeSecondDuration = 3 * time.Second

type Server struct {
	name       string
	httpServer *http.Server
	port       string
	timeout    time.Duration
	shutdownCH chan os.Signal
}

func New(name, port string, timeout time.Duration) *Server {
	shutdownCH := make(chan os.Signal, 1)
	signal.Notify(shutdownCH, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	return &Server{
		name:       name,
		port:       port,
		timeout:    timeout,
		shutdownCH: shutdownCH,
	}
}

func (s *Server) Start() {
	handler := s.registerRoutes()
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", s.port),
		Handler:           handler,
		WriteTimeout:      ThreeSecondDuration,
		ReadHeaderTimeout: ThreeSecondDuration,
	}
	s.httpServer = server

	done := make(chan struct{}, 1)
	go s.shutdownHook(done)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		os.Exit(1)
	}

	if err != nil {
		s.shutdownCH <- syscall.SIGINT
	}

	<-done
}
