package server

import (
	"log"
	"log/slog"
	"regexp"
	"strings"

	"github.com/0sokrat0/GoApiYA/orchestrator/internal/expr"
	genid "github.com/0sokrat0/GoApiYA/orchestrator/pkg/GenID"
	"github.com/0sokrat0/GoApiYA/orchestrator/pkg/db"
	"github.com/gofiber/fiber/v2"
)

type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type CalculateRequest struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

func (s *Server) CreateExpression(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return c.Status(415).JSON(fiber.Map{"error": "Unsupported Content-Type. Ожидается application/json"})
	}

	var req ExpressionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(422).JSON(fiber.Map{"error": "Невалидный JSON в запросе"})
	}

	req.Expression = strings.TrimSpace(req.Expression)
	if req.Expression == "" {
		return c.Status(422).JSON(fiber.Map{"error": "Поле 'expression' не может быть пустым"})
	}

	validExpr := regexp.MustCompile(`^[0-9+\-*/\s()]+$`)
	if !validExpr.MatchString(req.Expression) {
		return c.Status(422).JSON(fiber.Map{"error": "Выражение содержит недопустимые символы"})
	}

	slog.Info("Получено выражение", "expression", req.Expression)

	if err := s.ProcessExpression(req.Expression, c); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	return nil
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
	for i, expr := range expressions {
		if expr.RootTaskID != "" {
			if rootTask, found := s.store.Task.GetByID(expr.RootTaskID); found {
				if rootTask.Status == db.StatusSuccess {
					expr.Result = rootTask.Result
					expr.Status = db.StatusSuccess
					s.store.Expression.Update(expr)
					expressions[i] = expr
				}
			}
		}
	}

	return c.Status(200).JSON(expressions)
}

func (s *Server) GetByID(c *fiber.Ctx) error {
	exprID := c.Params("id")
	expr, ok := s.store.Expression.GetByID(exprID)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "Expression not found"})
	}

	if expr.RootTaskID != "" {
		rootTask, found := s.store.Task.GetByID(expr.RootTaskID)
		if found && rootTask.Status == db.StatusSuccess {
			expr.Result = rootTask.Result
			expr.Status = db.StatusSuccess
			s.store.Expression.Update(*expr)
		}
	}

	return c.Status(200).JSON(expr)
}

func (s *Server) GetTasks(c *fiber.Ctx) error {
	task, found, err := s.store.Task.GetNextTask()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if !found {
		return c.Status(404).JSON(fiber.Map{"error": "No tasks available"})
	}

	return c.Status(200).JSON(task)
}

func (s *Server) GetTasksAll(c *fiber.Ctx) error {
	tasks, err := s.store.Task.List()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(tasks)
}

func (s *Server) UpdateTasks(c *fiber.Ctx) error {
	var req CalculateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(422).JSON(fiber.Map{"error": "Invalid request"})
	}

	task, err := s.store.Task.Update(req.ID, req.Result)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	s.updateParentTasks(req.ID, req.Result)

	return c.Status(200).JSON(task)
}

func (s *Server) updateParentTasks(taskID string, result float64) {
	tasks, err := s.store.Task.List()
	if err != nil {

		log.Printf("Error listing tasks: %v", err)
		return
	}
	for _, parent := range tasks {
		updated := false
		if parent.LeftTaskID == taskID {
			parent.Arg1 = result
			parent.LeftTaskID = ""
			updated = true
		}
		if parent.RightTaskID == taskID {
			parent.Arg2 = result
			parent.RightTaskID = ""
			updated = true
		}
		if updated {
			if parent.LeftTaskID == "" && parent.RightTaskID == "" {
				parent.Status = db.StatusReady
			}
			s.store.Task.Replace(parent)
		}
	}
}

// func (s *Server) updateExpressionResult(exprID, rootTaskID string) {
// 	rootTask, ok := s.store.Task.GetByID(rootTaskID)
// 	if !ok {
// 		slog.Error("Корневая задача не найдена", "rootTaskID", rootTaskID)
// 		return
// 	}
// 	if rootTask.Status == db.StatusSuccess {
// 		if err := s.store.Expression.UpdateResult(exprID, rootTask.Result); err != nil {
// 			slog.Error("Не удалось обновить результат выражения", "exprID", exprID, "error", err)
// 		} else {
// 			slog.Info("Результат выражения обновлён", "exprID", exprID, "result", rootTask.Result)
// 		}
// 	}
// }
