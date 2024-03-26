package http

import (
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	port string
}

func NewServer(
	port string,
) Server {
	return Server{
		port: port,
	}
}

func (s *Server) Run() {
	app := fiber.New()
	routes := newRouter()
	routes.serve(app)
	_ = app.Listen(":" + s.port)
}
