package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {}

func NewProductController() *ProductController {
	return &ProductController{}
}

// Exemple Function
// @Summary      Get Server Info
// @Description  Get server info and dependencies status and uptime of server and more
// @Tags         ProductController
// @Produce      json
// @Success      200				 {object}  services.ServerInfoResponse
// @Router       /api/v1/ProductController/info [get]
func (c *ProductController) Example(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello from ProductController")
}
