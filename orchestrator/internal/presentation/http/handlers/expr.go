package handlers

import (
	"errors"
	"regexp"
	"strings"

	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/app/expr"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/presentation/http/dto"
	grpcClient "github.com/0sokrat0/GoApiYA/orchestrator/pkg/gen/api/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var ErrNotFound = errors.New("not found")

type Handlers struct {
	fiber *fiber.App
	cfg   *config.Config
	uc    expr.CalcOrchUsecase
	log   *zap.Logger
	auth  grpcClient.AuthClient
}

func NewHandlers(fiber *fiber.App,
	cfg *config.Config,
	uc expr.CalcOrchUsecase,
	log *zap.Logger,
	auth grpcClient.AuthClient,
) *Handlers {
	return &Handlers{
		fiber: fiber,
		cfg:   cfg,
		uc:    uc,
		log:   log,
		auth:  auth,
	}
}

func (h *Handlers) CreateExpression(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return c.Status(415).JSON(fiber.Map{"error": "Unsupported Content-Type. Ожидается application/json"})
	}

	ownerID, err := h.getOwnerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": err.Error()})
	}

	var req dto.ExpressionRequest
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

	h.log.Info("Получено выражение", zap.String("expression", req.Expression))

	if err := h.ProcessExpression(req.Expression, ownerID, c); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	return nil
}

func (h *Handlers) getOwnerID(c *fiber.Ctx) (string, error) {
	raw := c.Locals("jwt")
	token, ok := raw.(*jwt.Token)
	if !ok || !token.Valid {
		return "", errors.New("invalid JWT token")
	}
	claims, ok := token.Claims.(*dto.CustomClaims)
	if !ok {
		return "", errors.New("cannot parse custom JWT claims")
	}
	if claims.Subject == "" {
		return "", errors.New("sub claim is missing")
	}
	return claims.Subject, nil
}

func (h *Handlers) ProcessExpression(expression, ownerID string, c *fiber.Ctx) error {

	e, err := h.uc.CreateExp(expression, ownerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"id": e.ID})
}

func (h *Handlers) GetListExpressions(c *fiber.Ctx) error {
	ownerID, err := h.getOwnerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": err.Error()})
	}

	expressions, err := h.uc.ListExpressions(ownerID)
	if err != nil {
		h.log.Error("ListExpressions failed", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": "Внутренняя ошибка"})
	}
	return c.JSON(expressions)
}

func (h *Handlers) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	expr, err := h.uc.GetExpression(id)
	switch {
	case errors.Is(err, ErrNotFound):
		return c.Status(404).JSON(fiber.Map{"error": "Expression not found"})
	case err != nil:
		h.log.Error("GetExpression failed", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": "Внутренняя ошибка"})
	}

	return c.JSON(expr)
}

func (h *Handlers) GetTasksAll(c *fiber.Ctx) error {
	tasks, err := h.uc.ListTasks()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(tasks)
}
