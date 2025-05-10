package expr_test

import (
	"testing"

	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/app/expr"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func getTestConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			TIME_ADDITION_MS:       1000,
			TIME_SUBTRACTION_MS:    1000,
			TIME_MULTIPLICATION_MS: 1000,
			TIME_DIVISION_MS:       1000,
		},
	}
}

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task entity.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(id string) (entity.Task, bool, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Task), args.Bool(1), args.Error(2)
}
func (m *MockTaskRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskRepository) ListResponses() ([]entity.Task, error) {
	args := m.Called()
	return args.Get(0).([]entity.Task), args.Error(1)
}

// Для MockExpressionRepository
func (m *MockExpressionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockExpressionRepository) UpdateResult(exprID string, result float64) error {
	args := m.Called(exprID, result)
	return args.Error(0)
}

func (m *MockTaskRepository) List() ([]entity.Task, error) {
	args := m.Called()
	return args.Get(0).([]entity.Task), args.Error(1)
}

func (m *MockTaskRepository) Replace(task entity.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetNextTask() (entity.Task, bool, error) {
	args := m.Called()
	return args.Get(0).(entity.Task), args.Bool(1), args.Error(2)
}

func (m *MockTaskRepository) Update(id string, result float64) (*entity.Task, error) {
	args := m.Called(id, result)
	return args.Get(0).(*entity.Task), args.Error(1)
}

// MockExpressionRepository реализация
type MockExpressionRepository struct {
	mock.Mock
}

func (m *MockExpressionRepository) Create(expr entity.Expression) error {
	args := m.Called(expr)
	return args.Error(0)
}

func (m *MockExpressionRepository) GetByID(id string) (*entity.Expression, bool) {
	args := m.Called(id)
	return args.Get(0).(*entity.Expression), args.Bool(1)
}

func (m *MockExpressionRepository) Update(expr entity.Expression) error {
	args := m.Called(expr)
	return args.Error(0)
}

func (m *MockExpressionRepository) List(ownerID string) ([]entity.Expression, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]entity.Expression), args.Error(1)
}

func TestCalcOrch_CreateExp(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockExpRepo := new(MockExpressionRepository)
	logger := zap.NewNop()

	// Исправляем ожидание количества создаваемых задач
	// Для выражения "2+2*2" должно быть 2 задачи: умножение и сложение
	mockTaskRepo.On("Create", mock.Anything).Return(nil).Twice()
	mockExpRepo.On("Create", mock.Anything).Return(nil)
	mockExpRepo.On("Update", mock.Anything).Return(nil)

	service := expr.NewCalcOrch(
		mockTaskRepo,
		mockExpRepo,
		getTestConfig(),
		logger,
	)

	expression, err := service.CreateExp("2+2*2", "user1")

	assert.NoError(t, err)
	assert.NotEmpty(t, expression.ID)
	mockExpRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}
func TestCalcOrch_GetExpression(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockExpRepo := new(MockExpressionRepository)
	logger := zap.NewNop()

	exprID := uuid.NewString()
	taskID := uuid.NewString()

	testExpr := &entity.Expression{
		ID:         exprID,
		RootTaskID: taskID,
		Status:     entity.StatusPending,
	}

	mockExpRepo.On("GetByID", exprID).Return(testExpr, true)
	mockTaskRepo.On("GetByID", taskID).Return(
		entity.Task{ID: taskID, Status: entity.StatusSuccess, Result: 8.0},
		true,
		nil,
	)
	mockExpRepo.On("Update", mock.Anything).Return(nil)

	service := expr.NewCalcOrch(
		mockTaskRepo,
		mockExpRepo,
		getTestConfig(),
		logger,
	)

	result, err := service.GetExpression(exprID)

	assert.NoError(t, err)
	assert.Equal(t, entity.StatusSuccess, result.Status)
	assert.Equal(t, 8.0, result.Result)
}
func TestCalcOrch_ProcessTaskResult(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockExpRepo := new(MockExpressionRepository)
	logger := zap.NewNop()

	taskID := uuid.NewString()
	parentTaskID := uuid.NewString()

	mockTaskRepo.On("List").Return([]entity.Task{
		{
			ID:          parentTaskID,
			LeftTaskID:  taskID,
			RightTaskID: "other",
			Status:      entity.StatusPending,
		},
	}, nil)

	mockTaskRepo.On("Replace", mock.MatchedBy(func(t entity.Task) bool {
		return t.ID == parentTaskID && t.Arg1 == 5.0 && t.LeftTaskID == ""
	})).Return(nil)

	service := expr.NewCalcOrch(
		mockTaskRepo,
		mockExpRepo,
		getTestConfig(),
		logger,
	)

	// Вызов метода
	err := service.ProcessTaskResult(taskID, 5.0)

	// Проверки
	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
}

func TestCalcOrch_GetNextTask(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockExpRepo := new(MockExpressionRepository)
	logger := zap.NewNop()

	expectedTask := entity.Task{
		ID:        uuid.NewString(),
		Status:    entity.StatusReady,
		Operation: "+",
		Arg1:      2,
		Arg2:      3,
	}

	// Настройка моков
	mockTaskRepo.On("GetNextTask").Return(expectedTask, true, nil)

	service := expr.NewCalcOrch(
		mockTaskRepo,
		mockExpRepo,
		getTestConfig(),
		logger,
	)

	// Вызов метода
	task, found, err := service.GetNextTask()

	// Проверки
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, expectedTask, task)
}

func TestCalcOrch_UpdateTaskResult(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockExpRepo := new(MockExpressionRepository)
	logger := zap.NewNop()

	taskID := uuid.NewString()
	updatedTask := &entity.Task{
		ID:     taskID,
		Result: 10.0,
		Status: entity.StatusSuccess,
	}

	// Настройка моков
	mockTaskRepo.On("Update", taskID, 10.0).Return(updatedTask, nil)

	service := expr.NewCalcOrch(
		mockTaskRepo,
		mockExpRepo,
		getTestConfig(),
		logger,
	)

	// Вызов метода
	result, err := service.UpdateTaskResult(taskID, 10.0)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, updatedTask, result)
	mockTaskRepo.AssertExpectations(t)
}

func TestCalcOrch_ListExpressions(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockExpRepo := new(MockExpressionRepository)
	logger := zap.NewNop()

	ownerID := "user1"
	exprs := []entity.Expression{
		{
			ID:         "1",
			OwnerID:    ownerID,
			Status:     entity.StatusPending,
			RootTaskID: "task1",
		},
		{
			ID:         "2",
			OwnerID:    ownerID,
			Status:     entity.StatusSuccess,
			RootTaskID: "task2",
		},
	}

	// Настройка моков
	mockExpRepo.On("List", ownerID).Return(exprs, nil)

	// Настройка для task1
	mockTaskRepo.On("GetByID", "task1").Return(
		entity.Task{Status: entity.StatusSuccess, Result: 10.0},
		true,
		nil,
	)

	// Настройка для task2
	mockTaskRepo.On("GetByID", "task2").Return(
		entity.Task{Status: entity.StatusSuccess, Result: 20.0},
		true,
		nil,
	)

	// Ожидаем два вызова Update с любым выражением
	mockExpRepo.On("Update", mock.Anything).Return(nil).Twice()

	service := expr.NewCalcOrch(
		mockTaskRepo,
		mockExpRepo,
		getTestConfig(),
		logger,
	)

	result, err := service.ListExpressions(ownerID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// Проверяем обновление статусов
	for _, expr := range result {
		assert.Equal(t, entity.StatusSuccess, expr.Status)
	}

	mockExpRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}
