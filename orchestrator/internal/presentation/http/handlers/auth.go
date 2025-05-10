package handlers

import (
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/presentation/http/dto"
	authpb "github.com/0sokrat0/GoApiYA/orchestrator/pkg/gen/api/auth"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handlers) Auth(c *fiber.Ctx) error {
	raw := c.Body()
	if raw == nil {
		return status.Error(codes.Internal, "empty body")
	}

	h.log.Info("Auth: raw body", zap.ByteString("body", raw))
	var req dto.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Auth: invalid JSON", zap.Error(err))
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": "Невалидный JSON в запросе"})
	}
	if req.Login == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Поля 'login' и 'password' обязательны"})
	}

	grpcResp, err := h.auth.Login(c.Context(), &authpb.LoginRequest{
		Email:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		st, _ := status.FromError(err)
		h.log.Info("Auth: grpc.Login failed", zap.String("code", st.Code().String()), zap.String("desc", st.Message()))

		switch st.Code() {
		case codes.Unauthenticated, codes.PermissionDenied:
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "Неправильный логин или пароль"})
		case codes.InvalidArgument:
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": st.Message()})
		default:
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "Внутренняя ошибка аутентификации"})
		}
	}

	accessToken := grpcResp.Access
	refreshToken := grpcResp.Refresh

	accessCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   false,  // Только для HTTPS в production
		SameSite: "None", // Или "Strict" вместо "None"
		Path:     "/",
		MaxAge:   3600,
	}

	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Path:     "/refresh",
		MaxAge:   7 * 24 * 3600,
	}

	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Authenticated",
	})
}

func (h *Handlers) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Register: invalid JSON", zap.Error(err))
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": "Невалидный JSON в запросе"})
	}
	if req.Login == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Поля 'login' и 'password' обязательны"})
	}

	ctx := c.Context()
	_, err := h.auth.Register(ctx, &authpb.RegisterRequest{
		Email:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		h.log.Error("Register: grpc.Register failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Не удалось зарегистрировать пользователя"})
	}

	return c.Status(fiber.StatusCreated).
		JSON(fiber.Map{"message": "Регистрация успешна"})
}
