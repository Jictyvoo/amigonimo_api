package web

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/golang-jwt/jwt/v5"

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

func NewServer(conf config.Config, serverOptions ...ServerOption) *Server {
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
		fuego.WithMaxBodySize(2 * 1024 * 1024),
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
			jwtware.Config[jwt.MapClaims]{
				SigningKey: jwtware.SigningKey{
					JWTAlg: "PS256",
					Key:    conf.Runtime.AuthSecretKey,
				},
			},
		),
	}
	for _, opt := range serverOptions {
		opt(server)
	}

	_ = server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() error {
	// Firstly do with public APIs
	if err := SetupRoutes(s.Server, s.publicRouters...); err != nil {
		return err
	}

	// Add the authMiddleware on protected routers
	for _, router := range s.protectedRouters {
		router.AddMiddleware(s.authMiddleware)
	}

	err := SetupRoutes(s.Server, s.protectedRouters...)
	return err
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

func (s *Server) gracefulShutdown(sigint chan os.Signal, notifyComplete func()) {
	defer notifyComplete()
	// Start signal detection in a goroutine
	receivedSignal := <-sigint
	slog.With("signal", receivedSignal.String()).Warn("Server shutdown init")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Server shutdown complete")
}
