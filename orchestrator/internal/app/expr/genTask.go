package expr

import (
	"fmt"
	"strconv"

	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/entity"
	"github.com/google/uuid"
)

func GenerateTasks(expr Expr, cfg *config.Config) (string, []entity.Task) {
	var tasks []entity.Task

	var traverse func(e Expr) string
	traverse = func(e Expr) string {
		switch node := e.(type) {
		case Number:
			return strconv.FormatFloat(node.Value, 'f', -1, 64)
		case BinaryExpr:
			leftResult := traverse(node.Left)
			rightResult := traverse(node.Right)

			taskID := uuid.NewString()

			var task entity.Task
			task.ID = taskID
			task.Operation = operatorToString(node.Operator)
			task.OperationTime = getOperationTime(task.Operation, cfg)

			leftVal, leftErr := strconv.ParseFloat(leftResult, 64)
			if leftErr == nil {
				task.Arg1 = leftVal
			} else {
				task.LeftTaskID = leftResult
			}

			rightVal, rightErr := strconv.ParseFloat(rightResult, 64)
			if rightErr == nil {
				task.Arg2 = rightVal
			} else {
				task.RightTaskID = rightResult
			}

			task.Result = 0

			if leftErr == nil && rightErr == nil {
				task.Status = entity.StatusReady
			} else {
				task.Status = entity.StatusWaiting
			}

			tasks = append(tasks, task)
			return taskID
		default:
			panic(fmt.Sprintf("неподдерживаемый тип узла: %T", e))
		}
	}

	rootResult := traverse(expr)
	return rootResult, tasks
}

func operatorToString(op TokenType) string {
	switch op {
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenMul:
		return "*"
	case TokenDiv:
		return "/"
	default:
		return "unknown"
	}
}

func getOperationTime(op string, cfg *config.Config) int64 {
	switch op {
	case "+":
		return cfg.App.TIME_ADDITION_MS
	case "-":
		return cfg.App.TIME_SUBTRACTION_MS
	case "*":
		return cfg.App.TIME_MULTIPLICATION_MS
	case "/":
		return cfg.App.TIME_DIVISION_MS
	default:
		return 1000
	}
}
