package router

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/megorka/goproject/authorization/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Config struct {
	Host         string `yaml:"HTTP_HOST" env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port         string `yaml:"HTTP_PORT" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout  int    `yaml:"HTTP_READ_TIMEOUT" env:"HTTP_READ_TIMEOUT" env-default:"10"` // в секундах
	WriteTimeout int    `yaml:"HTTP_WRITE_TIMEOUT" env:"HTTP_WRITE_TIMEOUT" env-default:"30"`
	IdleTimeout  int    `yaml:"HTTP_IDLE_TIMEOUT" env:"HTTP_IDLE_TIMEOUT" env-default:"60"`
}

type Router struct {
	server  *http.Server
	config  Config
	Router  *mux.Router
	Handler *Handler
}

func NewRouter(cfg Config, h *Handler) *Router {
	r := mux.NewRouter()
	r.Use(MiddleWare)
	r.HandleFunc("/api/v1/auth/signup", h.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", h.Login).Methods("POST")
	r.HandleFunc("/api/v1/auth/google", h.GoogleLogin).Methods("GET")
	r.HandleFunc("/api/v1/auth/google/callback", h.GoogleCallback).Methods("GET")

	return &Router{
		config:  cfg,
		Router:  r,
		Handler: h,
	}
}

func (r *Router) Run(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%s", r.config.Host, r.config.Port)
	r.server = &http.Server{
		Addr:         addr,
		Handler:      r.Router,
		ReadTimeout:  time.Duration(r.config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(r.config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(r.config.IdleTimeout) * time.Second,
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Starting server", zap.String("address", addr))
	return r.server.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	if r.server == nil {
		return nil
	}
	return r.server.Shutdown(ctx)
}
