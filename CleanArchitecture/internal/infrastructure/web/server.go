package web

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server представляет HTTP сервер
type Server struct {
	httpServer *http.Server
	handler    http.Handler
}

// Config конфигурация сервера
type Config struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// NewServer создает новый HTTP сервер
func NewServer(config Config, handler http.Handler) *Server {
	// Устанавливаем значения по умолчанию если не заданы
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 15 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 15 * time.Second
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = 60 * time.Second
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		handler:    handler,
	}
}

// Start запускает HTTP сервер
func (s *Server) Start() error {
	fmt.Printf("Starting HTTP server on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully останавливает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}

// ПРИНЦИПЫ:
// 1. Простая обертка над стандартным http.Server
// 2. Конфигурируемые таймауты для безопасности
// 3. Graceful shutdown для корректного завершения
// 4. Логирование состояния сервера