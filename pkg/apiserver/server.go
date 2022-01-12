package apiserver

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/handlers"
	"github.com/teal-seagull/lyre-be-v4/pkg/config"
	"github.com/teal-seagull/lyre-be-v4/pkg/log"
)

type APIHTTPServer struct {
	apiVersion string
	srv        *http.Server
	logger     *zap.Logger
}

// NewAPIHTTPServer API server constructor
func NewAPIHTTPServer(cfg config.Server) (*APIHTTPServer, error) {
	s := &APIHTTPServer{
		apiVersion: cfg.APIVersion,
		logger:     log.TheLogger().With(zap.String("component", "APIHTTPServer")),
	}

	s.srv = &http.Server{
		Handler: handlers.SetupRouter(cfg),
	}

	if cfg.Address != "" {
		s.srv.Addr = cfg.Address
	}

	return s, nil
}

// Run starts HTTP server, ctx is used for server shutdown in case if ctx is closed
func (s *APIHTTPServer) Run(ctx context.Context) {
	loggerWithField := s.logger.With(zap.String("method", "Run"))

	go func() {
		for {
			<-ctx.Done()
			shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
			_ = s.srv.Shutdown(shutdownCtx)
			cancelFn()
			return
		}
	}()

	loggerWithField.With(zap.String("address", s.srv.Addr), zap.String("apiVersion", s.apiVersion)).
		Info("start server")

	if err := s.srv.ListenAndServe(); err != nil {
		loggerWithField.Warn("http server finished with error", zap.Error(err))
	}
}
