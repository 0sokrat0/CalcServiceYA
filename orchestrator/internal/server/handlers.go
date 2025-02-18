package server

import "github.com/gofiber/fiber/v2"

func (s *Server) Calculate(c *fiber.Ctx) error {
	return c.SendString("Hello, World !")
}
