package http

import "github.com/gofiber/fiber/v2"

func (s *Server) SetupRoutes(jwtMW fiber.Handler) {

	api := s.fiber.Group("api")
	api.Post("/auth", s.handlers.Auth)
	api.Post("/register", s.handlers.Register)

	v1 := api.Group("/v1", jwtMW)
	v1.Post("/calculate", s.handlers.CreateExpression)
	v1.Get("/expressions", s.handlers.GetListExpressions)
	v1.Get("/expressions/:id", s.handlers.GetByID)

	internal := s.fiber.Group("/internal")
	internal.Get("/tasks/", s.handlers.GetTasksAll)

	s.fiber.Static("/", "./public")

}
