package http

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) swaggerRouter() {
	r.fiberRouter.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/doc.json",
	}))

	r.fiberRouter.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(swagger.New())
	})
}
