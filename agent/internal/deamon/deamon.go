package deamon

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/0sokrat0/GoApiYA/agent/config"
)

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
}

func CalcPool(tasks <-chan Task, cfg *config.Config) {
	workers := cfg.App.COMPUTING_POWER
	slog.Info("Запуск пула вычислителей", "workers", workers)
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			for task := range tasks {
				slog.Info("Рабочий начал обработку задачи", "workerID", workerID, "taskID", task.ID)
				calc(task, cfg)
				slog.Info("Рабочий завершил обработку задачи", "workerID", workerID, "taskID", task.ID)
			}
		}(i)
	}
}

func calc(task Task, cfg *config.Config) {
	slog.Info("Начало вычислений", "taskID", task.ID, "Операция", task.Operation, "Arg1", task.Arg1, "Arg2", task.Arg2)

	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	slog.Info("Задержка окончена", "taskID", task.ID, "Операция", task.Operation)

	var result float64
	var opErr error

	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
		slog.Info("Выполнено сложение", "taskID", task.ID, "результат", result)
	case "-":
		result = task.Arg1 - task.Arg2
		slog.Info("Выполнено вычитание", "taskID", task.ID, "результат", result)
	case "*":
		result = task.Arg1 * task.Arg2
		slog.Info("Выполнено умножение", "taskID", task.ID, "результат", result)
	case "/":
		if task.Arg2 == 0 {
			opErr = fmt.Errorf("деление на ноль")
		} else {
			result = task.Arg1 / task.Arg2
			slog.Info("Выполнено деление", "taskID", task.ID, "результат", result)
		}
	default:
		opErr = fmt.Errorf("неподдерживаемая операция: %s", task.Operation)
		slog.Error("Обнаружена неподдерживаемая операция", "taskID", task.ID, "Операция", task.Operation)
	}

	if opErr != nil {
		slog.Error("Ошибка операции", "taskID", task.ID, "ошибка", opErr)
		return
	}

	slog.Info("Отправка вычисленного результата", "taskID", task.ID, "результат", result)
	if err := ResultTask(*cfg, result, task.ID); err != nil {
		slog.Error("Ошибка отправки результата", "taskID", task.ID, "ошибка", err)
		return
	}

	slog.Info(fmt.Sprintf("Задача %s завершена с результатом: %.2f", task.ID, result))
}

func ResultTask(cfg config.Config, result float64, id string) error {
	url := fmt.Sprintf("http://%s:%s/internal/task", cfg.Server.Host, cfg.Server.Port)
	payload := fmt.Sprintf(`{"id": "%s", "result": %.2f}`, id, result)
	slog.Info("Отправка результата на сервер", "URL", url, "payload", payload)

	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		slog.Error("Ошибка http.Post", "taskID", id, "ошибка", err)
		return fmt.Errorf("http.Post error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Сервер вернул не OK", "taskID", id, "status", resp.Status)
		return fmt.Errorf("failed to send result, status: %s", resp.Status)
	}

	slog.Info("Результат успешно отправлен", "taskID", id)
	return nil
}

func GetTask(cfg config.Config, tasksChan chan<- Task) {
	for {
		url := fmt.Sprintf("http://%s:%s/internal/task", cfg.Server.Host, cfg.Server.Port)
		resp, err := http.Get(url)
		if err != nil {
			slog.Error("Ошибка получения задачи", "ошибка", err)
			time.Sleep(10 * time.Second)
			continue
		}

		var task Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			slog.Error("Ошибка декодирования задачи", "ошибка", err)
			resp.Body.Close()
			time.Sleep(10 * time.Second)
			continue
		}
		resp.Body.Close()

		if task.ID == "" {
			time.Sleep(1 * time.Second)
			continue
		}

		slog.Info("Задача получена", "taskID", task.ID)
		tasksChan <- task

		time.Sleep(1 * time.Second)
	}
}
