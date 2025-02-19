package db

import "sync"

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	LeftTaskID    string  `json:"left_task_id"`
	RightTaskID   string  `json:"right_task_id"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
	Result        float64 `json:"result"`
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
