package server

import (
	"CalcYA/orchestrator/internal/calc"
	genid "CalcYA/orchestrator/pkg/GenID"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CalculateRequest struct {
	Expression string `json:"expression"`
}

func (s *Server) CreateExpression(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")
	var req CalculateRequest
	switch {
	case contentType == "application/json":
		if err := c.BodyParser(&req); err != nil {
			return c.Status(422).JSON(fiber.Map{"error": "invalid request"})
		}
		req.Expression = strings.TrimSpace(req.Expression)
		if req.Expression == "" {
			return c.Status(422).JSON(fiber.Map{"error": "Поле 'expression' не может быть пустым"})
		}
		return s.ProcessExpression(req.Expression, c)
	}
	return s.ProcessExpression(req.Expression, c)
}

func (s *Server) ProcessExpression(expr string, c *fiber.Ctx) error {
	id := genid.GenerateID()
	_, err := calc.CreateExp(s.store, expr, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"id": id})
}

func (s *Server) GetListExpressions(c *fiber.Ctx) error {
	expressions, err := s.store.List()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(expressions)
}

func (s *Server) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	expr, ok := s.store.GetByID(id)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "Expression not found"})
	}

	if err := c.Status(200).JSON(expr); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return nil
}
