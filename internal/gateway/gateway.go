package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"taskhub/config"
	"taskhub/internal/app"
	"taskhub/internal/handler"
	"taskhub/pkg/logger"
	"taskhub/pkg/middleware"
	"taskhub/pkg/nats"

	"go.uber.org/fx"
)

var GatewayModule = fx.Module(
	"gateway",
	fx.Provide(NewGateway),
)

type Gateway struct {
	config       *config.Config
	natsConn     *nats.Nats
	httpServer   *http.Server
	logger       *logger.Logger
	authHandler  *handler.AuthHandler
	taskHandler  *handler.TaskHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewGateway(
	config *config.Config,
	logger *logger.Logger,
	authService *app.AuthService,
	taskService *app.TaskService,
) *Gateway {
	return &Gateway{
		config:       config,
		natsConn:     nats.NewNats(config, logger),
		logger:       logger,
		authHandler:  handler.NewAuthHandler(authService),
		taskHandler:  handler.NewTaskHandler(taskService),
		authMiddleware: middleware.NewAuthMiddleware(authService),
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (g *Gateway) Start() error {
	if g.config == nil {
		return errors.New("config is nil")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthCheck)

	mux.HandleFunc("/api/auth/register", g.authHandler.Register)
	mux.HandleFunc("/api/auth/login", g.authHandler.Login)
	mux.HandleFunc("/api/auth/refresh", g.authHandler.RefreshToken)

	mux.Handle("/api/tasks", g.authMiddleware.Authenticate(http.HandlerFunc(g.handleTasks)))
	mux.Handle("/api/tasks/", g.authMiddleware.Authenticate(http.HandlerFunc(g.handleTaskByID)))

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

func (g *Gateway) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		g.taskHandler.List(w, r)
	case http.MethodPost:
		g.taskHandler.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (g *Gateway) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")

	if strings.HasSuffix(path, "/complete") {
		g.taskHandler.Complete(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		g.taskHandler.Get(w, r)
	case http.MethodPut:
		g.taskHandler.Update(w, r)
	case http.MethodDelete:
		g.taskHandler.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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
