package server

import (
	"context"

	"go.uber.org/zap"
)

// nolint
func (s *Server) shutdownHook(done chan<- struct{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	<-s.shutdownCH

	baseCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer func() {
		cancel()
	}()

	s.httpServer.SetKeepAlivesEnabled(false)
	err := s.httpServer.Shutdown(baseCtx)
	if err != nil {
		sugar.Errorf("error while shutting down the server gracefully, got error: %v", err)
	}

	close(done)
}
