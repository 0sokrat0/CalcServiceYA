package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	// Импортируйте ваш сервер и связанные пакеты. Пути могут отличаться.
	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/server"
	"github.com/0sokrat0/GoApiYA/orchestrator/pkg/db"
)

// setupTestServer создаёт тестовое Fiber-приложение с зарегистрированными endpoint-ами.
func setupTestServer() (*fiber.App, *server.Server) {
	app := fiber.New()

	cfg := &config.Config{}

	expressionStore := db.NewExpressionStore()
	taskStore := db.NewTaskStore()

	// Собираем хранилище
	store := &db.Stores{
		Expression: expressionStore,
		Task:       taskStore,
	}

	srv, err := server.NewServer(cfg, store)
	if err != nil {
		panic(err)
	}
	app.Post("/api/v1/calculate", srv.CreateExpression)
	app.Get("/api/v1/expressions", srv.GetListExpressions)
	app.Get("/api/v1/expressions/:id", srv.GetByID)
	app.Get("/internal/task", srv.GetTasks)
	app.Post("/internal/task", srv.UpdateTasks)

	return app, srv
}

func TestCreateExpressionValid(t *testing.T) {
	app, _ := setupTestServer()

	payload := `{"expression": "2+2*2"}`
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result["id"])
}

func TestCreateExpressionEmpty(t *testing.T) {
	app, _ := setupTestServer()

	payload := `{"expression": "   "}`
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestCreateExpressionInvalidContentType(t *testing.T) {
	app, _ := setupTestServer()

	payload := `{"expression": "2+2*2"}`
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
}

func TestCreateExpressionInvalidCharacters(t *testing.T) {
	app, _ := setupTestServer()

	payload := `{"expression": "2+2*2a"}`
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestGetListExpressions(t *testing.T) {
	app, _ := setupTestServer()

	createPayload := `{"expression": "2+2*2"}`
	reqCreate := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(createPayload))
	reqCreate.Header.Set("Content-Type", "application/json")
	respCreate, err := app.Test(reqCreate)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode)

	req := httptest.NewRequest("GET", "/api/v1/expressions", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var expressions []struct {
		ID     string  `json:"id"`
		Status string  `json:"status"`
		Result float64 `json:"result"`
	}
	err = json.NewDecoder(resp.Body).Decode(&expressions)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(expressions), 1)
}

func TestGetExpressionByID_NotFound(t *testing.T) {
	app, _ := setupTestServer()

	req := httptest.NewRequest("GET", "/api/v1/expressions/non-existent-id", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetExpressionByID_Found(t *testing.T) {
	app, _ := setupTestServer()

	// Создаем выражение
	expression := "2+2*2"
	createPayload := `{"expression": "` + expression + `"}`
	reqCreate := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(createPayload))
	reqCreate.Header.Set("Content-Type", "application/json")
	respCreate, err := app.Test(reqCreate)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode)

	var result map[string]string
	err = json.NewDecoder(respCreate.Body).Decode(&result)
	assert.NoError(t, err)
	id := result["id"]
	assert.NotEmpty(t, id)

	req := httptest.NewRequest("GET", "/api/v1/expressions/"+id, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var exprResp struct {
		ID     string  `json:"id"`
		Status string  `json:"status"`
		Result float64 `json:"result"`
	}
	err = json.NewDecoder(resp.Body).Decode(&exprResp)
	assert.NoError(t, err)
	assert.Equal(t, id, exprResp.ID)
}

func TestGetTasks_NoTask(t *testing.T) {
	app, _ := setupTestServer()

	req := httptest.NewRequest("GET", "/internal/task", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUpdateTasks_InvalidRequest(t *testing.T) {
	app, _ := setupTestServer()

	payload := `{"id":123, "result": "not a number"}`
	req := httptest.NewRequest("POST", "/internal/task", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}
