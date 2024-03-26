package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	fiberRouter fiber.Router
}

func newRouter() Router {
	return Router{}
}

func (r *Router) serve(app *fiber.App) {
	r.fiberRouter = app.Group("/")
	api := app.Group("/api")
	v1 := api.Group("/v1")
	fmt.Println(v1)

	r.swaggerRouter()
}
