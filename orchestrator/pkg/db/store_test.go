package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressionStore_CreateAndGet(t *testing.T) {
	store := NewExpressionStore()

	expr := Expression{
		ID:         "expr-1",
		RootTaskID: "task-1",
		Status:     StatusPending,
		Result:     0,
	}

	// Создаем выражение
	store.Create(expr)

	// Получаем выражение по ID
	got, ok := store.GetByID("expr-1")
	assert.True(t, ok, "Expression should be found")
	assert.Equal(t, expr.ID, got.ID)
	assert.Equal(t, expr.Status, got.Status)
}

func TestExpressionStore_ListAndUpdateResult(t *testing.T) {
	store := NewExpressionStore()

	expr1 := Expression{ID: "expr-1", RootTaskID: "task-1", Status: StatusPending, Result: 0}
	expr2 := Expression{ID: "expr-2", RootTaskID: "task-2", Status: StatusPending, Result: 0}
	store.Create(expr1)
	store.Create(expr2)

	list, err := store.List()
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	err = store.UpdateResult("expr-1", 42)
	assert.NoError(t, err)

	got, ok := store.GetByID("expr-1")
	assert.True(t, ok)
	assert.Equal(t, 42.0, got.Result)
	assert.Equal(t, StatusSuccess, got.Status)
}

func TestTaskStore_CreateGetListAndUpdate(t *testing.T) {
	store := NewTaskStore()

	task := Task{
		ID:            "task-1",
		Arg1:          2,
		Arg2:          3,
		LeftTaskID:    "",
		RightTaskID:   "",
		Operation:     "+",
		OperationTime: 100,
		Status:        StatusReady,
	}

	err := store.Create(task)
	assert.NoError(t, err)

	got, ok := store.GetByID("task-1")
	assert.True(t, ok)
	assert.Equal(t, task.Operation, got.Operation)

	updatedTask, err := store.Update("task-1", 5)
	assert.NoError(t, err)
	assert.Equal(t, 5.0, updatedTask.Result)
	assert.Equal(t, StatusSuccess, updatedTask.Status)

	list, err := store.List()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestTaskStore_GetNextTask(t *testing.T) {
	store := NewTaskStore()

	task1 := Task{
		ID:            "task-1",
		Arg1:          1,
		Arg2:          2,
		Operation:     "*",
		OperationTime: 100,
		Status:        StatusRunning, // не готова
	}
	task2 := Task{
		ID:            "task-2",
		Arg1:          3,
		Arg2:          4,
		Operation:     "+",
		OperationTime: 150,
		Status:        StatusReady, // готова
	}
	_ = store.Create(task1)
	_ = store.Create(task2)

	resp, found, err := store.GetNextTask()
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "task-2", resp.ID)

	taskAfter, ok := store.GetByID("task-2")
	assert.True(t, ok)
	assert.Equal(t, StatusRunning, taskAfter.Status)
}
