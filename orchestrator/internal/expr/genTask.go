package expr

import (
	"CalcYA/orchestrator/config"
	genid "CalcYA/orchestrator/pkg/GenID"
	"CalcYA/orchestrator/pkg/db"
	"fmt"
	"strconv"
)

func GenerateTasks(expr Expr, cfg *config.Config) (string, []db.Task) {
	var tasks []db.Task

	// Рекурсивная функция для обхода AST
	var traverse func(e Expr) string
	traverse = func(e Expr) string {
		switch node := e.(type) {
		case Number:
			return strconv.FormatFloat(node.Value, 'f', -1, 64)
		case BinaryExpr:
			leftResult := traverse(node.Left)
			rightResult := traverse(node.Right)

			taskID := genid.GenerateIDTask()

			var task db.Task
			task.ID = taskID
			task.Operation = operatorToString(node.Operator)
			task.OperationTime = getOperationTime(task.Operation, cfg)

			if val, err := strconv.ParseFloat(leftResult, 64); err == nil {
				task.Arg1 = val
			} else {
				task.LeftTaskID = leftResult
			}
			if val, err := strconv.ParseFloat(rightResult, 64); err == nil {
				task.Arg2 = val
			} else {
				task.RightTaskID = rightResult
			}

			task.Result = 0
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
		return cfg.App.TIME_MULTIPLICATIONS_MS
	case "/":
		return cfg.App.TIME_DIVISIONS_MS
	default:
		return 1000
	}
}
