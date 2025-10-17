package gateway

import (
	"context"
	"fmt"
	"net/http"
	"taskhub/config"
	"taskhub/pkg/logger"
	"taskhub/pkg/nats"
)

type Gateway struct {
	config     *config.Config
	natsConn   *nats.Nats
	httpServer *http.Server
	logger     *logger.Logger
}

func NewGateway(config *config.Config, logger *logger.Logger) *Gateway {
	return &Gateway{
		natsConn: nats.NewNats(config, logger),
		logger:   logger,
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleEvent(w http.ResponseWriter, r *http.Request) {}

func (g *Gateway) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/:events", handleEvent)
	mux.HandleFunc("/health", healthCheck)

	g.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", g.config.Port),
		Handler:      mux,
		ReadTimeout:  g.config.ReadTimeout,
		WriteTimeout: g.config.WriteTimeout,
		IdleTimeout:  g.config.IdleTimeout,
	}

	g.logger.Info("Starting HTTP server on port %s", g.config.Port)

	return g.httpServer.ListenAndServe()
}

func (g *Gateway) Shutdown(ctx context.Context) error {
	g.logger.Info("Shutting down HTTP server")

	if g.natsConn != nil {
		g.natsConn.Close()
	}

	if g.httpServer != nil {
		return g.httpServer.Shutdown(ctx)
	}

	return nil
}
