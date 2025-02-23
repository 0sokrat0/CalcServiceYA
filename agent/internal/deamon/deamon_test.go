package deamon_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/0sokrat0/GoApiYA/agent/config"
	"github.com/0sokrat0/GoApiYA/agent/internal/deamon"
	"github.com/stretchr/testify/assert"
)


func FakeTaskServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			
			w.Header().Set("Content-Type", "application/json")
			
			if strings.Contains(r.URL.Path, "empty") {
				w.Write([]byte(`{"id": ""}`))
			} else {
				w.Write([]byte(`{"id": "task-1", "arg1": 2, "arg2": 2, "operation": "*", "operation_time": 100}`))
			}
		case http.MethodPost:
			
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	return httptest.NewServer(handler)
}

func TestGetTaskIntegration(t *testing.T) {
	
	cfg, err := config.LoadConfig("../../config")
	assert.NoError(t, err)

	ts := FakeTaskServer()
	defer ts.Close()

	
	address := strings.TrimPrefix(ts.URL, "http://")
	cfg.Server.Host = address
	cfg.Server.Port = "" 

	
	tasksChan := make(chan deamon.Task, 1)
	go deamon.GetTask(*cfg, tasksChan)

	
	select {
	case task := <-tasksChan:
		assert.Equal(t, "task-1", task.ID, "Получена неверная задача")
		assert.Equal(t, 2.0, task.Arg1)
		assert.Equal(t, 2.0, task.Arg2)
		assert.Equal(t, "*", task.Operation)
	case <-time.After(5 * time.Second):
		t.Fatal("Время ожидания получения задачи истекло")
	}
}

func TestResultTaskIntegration(t *testing.T) {
	cfg, err := config.LoadConfig("../../config")

	assert.NoError(t, err)

	ts := FakeTaskServer()
	defer ts.Close()

	address := strings.TrimPrefix(ts.URL, "http://")
	cfg.Server.Host = address
	cfg.Server.Port = ""

	
	err = deamon.ResultTask(*cfg, 12.34, "task-1")
	assert.NoError(t, err, "Ошибка отправки результата")
}

func TestCalcPoolIntegration(t *testing.T) {
	cfg, err := config.LoadConfig("../../config")

	assert.NoError(t, err)

	
	tasksChan := make(chan deamon.Task, 1)
	testTask := deamon.Task{
		ID:            "task-2",
		Arg1:          3,
		Arg2:          4,
		Operation:     "*",
		OperationTime: 100,
	}
	tasksChan <- testTask
	close(tasksChan)


	done := make(chan struct{})
	go func() {
		deamon.CalcPool(tasksChan, cfg)
		close(done)
	}()

	select {
	case <-done:
		
	case <-time.After(3 * time.Second):
		t.Fatal("CalcPool не завершился в отведенное время")
	}
}
