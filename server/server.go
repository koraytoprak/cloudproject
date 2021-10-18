package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	return &Server{
		log: opts.Log,
	}
}

func New(opts Options) *Server {
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()
	return &Server{
		address: address,
		mux:     mux,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	s.setupRoutes()

	fmt.Println("Starting on", s.address)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

func (s *Server) Stop() error {
	fmt.Println("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}

func (s *Server) Start() error {

	s.log.Info("Starting", zap.String("address", s.address))

}

func (s *Server) Stop() error {
	s.log.Info("Stopping")

}

type Server struct {
	address string
	mux     chi.Router
	server  *http.Server
}

type Options struct {
	Host string
	Port int
}

type Server struct {
	log *zap.Logger
}

type Options struct {
	Log *zap.Logger
}
