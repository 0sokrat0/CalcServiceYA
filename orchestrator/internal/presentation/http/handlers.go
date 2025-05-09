// package http

// import (
// 	"errors"
// 	"regexp"
// 	"strings"

// 	"github.com/0sokrat0/GoApiYA/orchestrator/internal/presentation/http/dto"
// 	"github.com/gofiber/fiber/v2"
// 	"go.uber.org/zap"
// )

// var ErrNotFound = errors.New("not found")

// func (s *Server) CreateExpression(c *fiber.Ctx) error {
// 	contentType := c.Get("Content-Type")
// 	if !strings.Contains(contentType, "application/json") {
// 		return c.Status(415).JSON(fiber.Map{"error": "Unsupported Content-Type. Ожидается application/json"})
// 	}

// 	var req dto.ExpressionRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(422).JSON(fiber.Map{"error": "Невалидный JSON в запросе"})
// 	}

// 	req.Expression = strings.TrimSpace(req.Expression)
// 	if req.Expression == "" {
// 		return c.Status(422).JSON(fiber.Map{"error": "Поле 'expression' не может быть пустым"})
// 	}

// 	validExpr := regexp.MustCompile(`^[0-9+\-*/\s()]+$`)
// 	if !validExpr.MatchString(req.Expression) {
// 		return c.Status(422).JSON(fiber.Map{"error": "Выражение содержит недопустимые символы"})
// 	}

// 	s.log.Info("Получено выражение", zap.String("expression", req.Expression))

// 	if err := s.ProcessExpression(req.Expression, c); err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
// 	}

// 	return nil
// }

// func (s *Server) ProcessExpression(expression string, c *fiber.Ctx) error {

// 	e, err := s.uc.CreateExp(expression)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 	}
// 	return c.Status(201).JSON(fiber.Map{"id": e.ID})
// }

// func (s *Server) GetListExpressions(c *fiber.Ctx) error {
// 	expressions, err := s.uc.ListExpressions()
// 	if err != nil {
// 		s.log.Error("ListExpressions failed", zap.Error(err))
// 		return c.Status(500).JSON(fiber.Map{"error": "Внутренняя ошибка"})
// 	}
// 	return c.JSON(expressions)
// }

// func (s *Server) GetByID(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	expr, err := s.uc.GetExpression(id)
// 	switch {
// 	case errors.Is(err, ErrNotFound):
// 		return c.Status(404).JSON(fiber.Map{"error": "Expression not found"})
// 	case err != nil:
// 		s.log.Error("GetExpression failed", zap.Error(err))
// 		return c.Status(500).JSON(fiber.Map{"error": "Внутренняя ошибка"})
// 	}

// 	return c.JSON(expr)
// }

// func (s *Server) GetTasksAll(c *fiber.Ctx) error {
// 	tasks, err := s.uc.ListTasks()

//		if err != nil {
//			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
//		}
//		return c.Status(200).JSON(tasks)
//	}
package http
