package web

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/api"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
	"github.com/jictyvoo/amigonimo_api/pkg/web/middlewares/jwtware"
)

type ServerOption func(*Server)

func WithPublicRouters(routers ...RouterContract) ServerOption {
	return func(s *Server) {
		s.publicRouters = routers
	}
}

func WithPrivateRouters(routers ...RouterContract) ServerOption {
	return func(s *Server) {
		s.protectedRouters = routers
	}
}

type Server struct {
	*fuego.Server

	authMiddleware   HttpMiddleware
	protectedRouters []RouterContract
	publicRouters    []RouterContract
}

func NewServer(
	conf config.Config,
	jwtPublicKey *rsa.PublicKey,
	serverOptions ...ServerOption,
) (*Server, error) {
	const (
		kiloBytes    = 1024
		megaBytes    = kiloBytes * 1024
		maxBodyBytes = 2 * megaBytes
	)

	options := [4]func(*fuego.Server){
		fuego.WithEngineOptions(
			fuego.WithRequestContentType("application/json"),
			fuego.WithOpenAPIConfig(
				fuego.OpenAPIConfig{
					Disabled:     !conf.IsDebug,
					JSONFilePath: "api/openapi.json",
					Info: &openapi3.Info{
						Title:       "OpenAPI",
						Description: api.Description(),
						Version:     "0.0.1",
					},
				},
			),
		),
		fuego.WithDisallowUnknownFields(false),
		fuego.WithMaxBodySize(maxBodyBytes),
	}

	optSlice := options[:2] // Transform into a slice to save bytes later
	if conf.Runtime.Port != 0 {
		optSlice = append(
			optSlice, fuego.WithAddr(conf.Runtime.Host+":"+strconv.Itoa(int(conf.Runtime.Port))),
		)
	}

	fuegoServer := fuego.NewServer(optSlice...)
	server := &Server{
		Server: fuegoServer,
		authMiddleware: jwtware.New(
			jwtware.MapClaimsConfig{
				SigningKey: jwtware.SigningKey{
					JWTAlg: "PS256",
					Key:    jwtPublicKey,
				},
			},
		),
	}
	for _, opt := range serverOptions {
		opt(server)
	}

	if jwtPublicKey == nil {
		return nil, fmt.Errorf("JWT secret key is required but not provided or invalid")
	}

	if err := server.setupRoutes(); err != nil {
		return nil, fmt.Errorf("failed to setup routes: %w", err)
	}
	return server, nil
}

func (s *Server) Run() error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGSEGV)

	var wg sync.WaitGroup
	wg.Add(1)
	go s.gracefulShutdown(sigint, wg.Done)

	// Run server on the main goroutine (blocking)
	err := s.Server.Run()
	wg.Wait()
	return err
}

func (s *Server) setupRoutes() error {
	// Firstly do with public APIs
	if err := SetupRoutes(s.Server, s.publicRouters...); err != nil {
		return err
	}

	// Add the authMiddleware on protected routers
	for _, router := range s.protectedRouters {
		if withExtender, ok := router.(RouterMiddlewareExtender); ok {
			withExtender.AddMiddleware(s.authMiddleware)
		}
	}

	err := SetupRoutes(s.Server, s.protectedRouters...)
	return err
}

func (s *Server) gracefulShutdown(sigint chan os.Signal, notifyComplete func()) {
	defer notifyComplete()
	// Start signal detection in a goroutine
	receivedSignal := <-sigint
	slog.With("signal", receivedSignal.String()).Warn("Server shutdown init")

	const shutdownTimeout = 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Server shutdown complete")
}
