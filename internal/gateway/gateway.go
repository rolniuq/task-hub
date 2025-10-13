package gateway

import (
	"fmt"
	"net/http"
	"taskhub/config"
	"taskhub/pkg/nats"
)

type Gateway struct {
	config     *config.Config
	natsConn   *nats.Nats
	httpServer *http.Server
}

func NewGateway(config *config.Config) *Gateway {
	return &Gateway{
		natsConn: nats.NewNats(config),
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
		Addr: fmt.Sprintf(":%s", g.config.Port),
	}

	return g.httpServer.ListenAndServe()
}

func (g *Gateway) Shutdown() error {
	return nil
}
