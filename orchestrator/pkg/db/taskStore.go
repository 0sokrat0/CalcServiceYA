package db

import (
	"fmt"
	"sync"
)

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	LeftTaskID    string  `json:"left_task_id"`
	RightTaskID   string  `json:"right_task_id"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
	Result        float64 `json:"result"`
	Status        Status  `json:"status"`
}

type TasksResponse struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
}

type TaskStore struct {
	tasks map[string]Task
	mu    sync.RWMutex
}

type TaskRepository interface {
	List() ([]Task, error)
	Delete(id string) error
	Create(task Task) error
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]Task),
	}
}

func (s *TaskStore) List() ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskStore) ListResponses() ([]TasksResponse, error) {
	tasks, err := s.List()
	if err != nil {
		return nil, err
	}
	responses := make([]TasksResponse, 0, len(tasks))
	for _, task := range tasks {
		response := TasksResponse{
			ID:            task.ID,
			Arg1:          task.Arg1,
			Arg2:          task.Arg2,
			Operation:     task.Operation,
			OperationTime: task.OperationTime,
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (s *TaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tasks, id)
	return nil
}

func (s *TaskStore) Create(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.ID] = task
	return nil
}

func (s *TaskStore) Update(id string, result float64) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}

	task.Result = result
	task.Status = "success"
	s.tasks[id] = task
	return &task, nil
}

func (s *TaskStore) Replace(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[task.ID]; !ok {
		return fmt.Errorf("task not found")
	}
	s.tasks[task.ID] = task
	return nil
}

func (s *TaskStore) GetNextTask() (TasksResponse, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, task := range s.tasks {
		if task.Status == StatusReady {
			// Переводим задачу в статус running
			task.Status = StatusRunning
			s.tasks[id] = task

			// Формируем ответ в виде TasksResponse
			response := TasksResponse{
				ID:            task.ID,
				Arg1:          task.Arg1,
				Arg2:          task.Arg2,
				Operation:     task.Operation,
				OperationTime: task.OperationTime,
			}
			return response, true, nil
		}
	}
	return TasksResponse{}, false, nil
}

func (s *TaskStore) GetByID(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}
