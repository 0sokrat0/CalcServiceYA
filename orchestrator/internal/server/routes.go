package server

func (s *Server) SetupRoutes() {

	api := s.fiber.Group("api")

	v1 := api.Group("/v1")
	v1.Post("/calculate", s.CreateExpression)
	v1.Get("/expressions", s.GetListExpressions)
	v1.Get("/expressions/:id", s.GetByID)

	internal := s.fiber.Group("/internal")
	internal.Get("/task", s.GetTasks)
	internal.Get("/tasks/", s.GetTasksAll)
	internal.Post("/task", s.UpdateTasks)

}
