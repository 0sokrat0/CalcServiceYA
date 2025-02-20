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

// FakeTaskServer эмулирует поведение сервера оркестратора для GET и POST запросов по пути /internal/task.
func FakeTaskServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Для GET запроса возвращаем задачу, если доступна.
			// В тестовом примере будем возвращать задачу с фиксированными параметрами.
			// Чтобы имитировать ситуацию "нет задачи", можно возвращать пустой JSON.
			// Здесь возвращаем корректную задачу:
			w.Header().Set("Content-Type", "application/json")
			// Если URL содержит "empty", вернем пустую задачу.
			if strings.Contains(r.URL.Path, "empty") {
				w.Write([]byte(`{"id": ""}`))
			} else {
				w.Write([]byte(`{"id": "task-1", "arg1": 2, "arg2": 2, "operation": "*", "operation_time": 100}`))
			}
		case http.MethodPost:
			// Для POST запроса возвращаем статус OK.
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	return httptest.NewServer(handler)
}

func TestGetTaskIntegration(t *testing.T) {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("../../config")
	assert.NoError(t, err)

	// Запускаем фейковый сервер оркестратора
	ts := FakeTaskServer()
	defer ts.Close()

	// Перенастраиваем адрес сервера в конфиге для тестов
	// Пример: ts.URL = "http://127.0.0.1:12345"
	// Извлекаем адрес без схемы, так как наш код формирует URL с http://
	address := strings.TrimPrefix(ts.URL, "http://")
	cfg.Server.Host = address
	cfg.Server.Port = "" // т.к. ts.URL уже содержит порт

	// Создаем канал для задач
	tasksChan := make(chan deamon.Task, 1)
	go deamon.GetTask(*cfg, tasksChan)

	// Ожидаем получения задачи
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

	// Вызываем ResultTask для задачи task-1
	err = deamon.ResultTask(*cfg, 12.34, "task-1")
	assert.NoError(t, err, "Ошибка отправки результата")
}

func TestCalcPoolIntegration(t *testing.T) {
	cfg, err := config.LoadConfig("../../config")

	assert.NoError(t, err)

	// Создаем канал для задач и помещаем тестовую задачу
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

	// Для тестирования CalcPool, запускаем его и замеряем, что задача была обработана.
	// В CalcPool функция calc отправляет результат через ResultTask, поэтому мы эмулируем это.
	// Здесь просто проверим, что CalcPool завершает обработку.
	done := make(chan struct{})
	go func() {
		deamon.CalcPool(tasksChan, cfg)
		close(done)
	}()

	select {
	case <-done:
		// Если канал tasksChan был закрыт, CalcPool завершится.
		// Так как у нас всего одна задача, функция завершится после обработки.
	case <-time.After(3 * time.Second):
		t.Fatal("CalcPool не завершился в отведенное время")
	}
}
