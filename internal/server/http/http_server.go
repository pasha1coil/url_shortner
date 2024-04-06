package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Logger      *zap.Logger
	Controllers []Controller
}

type Server struct {
	Logger      *zap.Logger
	Controllers []Controller
	app         *fiber.App
}

func NewServer(config ServerConfig) *Server {
	app := fiber.New()

	s := &Server{
		Logger:      config.Logger,
		Controllers: config.Controllers,
		app:         app,
	}

	s.registerRoutes()

	return s
}

func (s *Server) Start(addr string) error {
	if err := s.app.Listen(addr); err != nil {
		s.Logger.Error("Failed to start server", zap.Error(err))
		return err
	}
	return nil
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func (s *Server) registerRoutes() {
	for _, c := range s.Controllers {
		router := s.app.Group(c.Name())
		c.Register(router)
	}
}

type Controller interface {
	Register(router fiber.Router)
	Name() string
}

func (s *Server) ListRoutes() {
	fmt.Println("Registered routes:")
	for _, stack := range s.app.Stack() {
		for _, route := range stack {
			fmt.Printf("%s %s\n", route.Method, route.Path)
		}
	}
}
