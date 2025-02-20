package server

import (
	"orchestrator/internal/expr"
	genid "orchestrator/pkg/GenID"
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

func (s *Server) ProcessExpression(expression string, c *fiber.Ctx) error {
	id := genid.GenerateID()
	_, err := expr.CreateExp(s.store.Expression, s.store.Task, expression, id, s.cfg)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"id": id})
}

func (s *Server) GetListExpressions(c *fiber.Ctx) error {
	expressions, err := s.store.Expression.List()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(expressions)
}

func (s *Server) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	expr, ok := s.store.Expression.GetByID(id)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "Expression not found"})
	}

	if err := c.Status(200).JSON(expr); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return nil
}

func (s *Server) GetTasks(c *fiber.Ctx) error {
	tasks, err := s.store.Task.List()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(tasks)
}
