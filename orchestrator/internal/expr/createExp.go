package expr

import (
	"CalcYA/orchestrator/config"

	"CalcYA/orchestrator/pkg/db"
	"fmt"
)

func CreateExp(exprStore *db.ExpressionStore, taskStore *db.TaskStore, expression string, id string, cfg *config.Config) (db.Expression, error) {
	exprRecord := db.Expression{
		ID:     id,
		Status: db.StatusAccepted,
		Result: 0,
	}
	exprStore.Create(exprRecord)

	tokens := tokenize(expression)
	parser := NewParser(tokens)
	ast := parser.parseExpression(0)

	rootID, tasks := GenerateTasks(ast, cfg)
	fmt.Printf("Корневой таск ID: %s\n", rootID)

	for _, task := range tasks {
		if err := taskStore.Create(task); err != nil {
			return exprRecord, err
		}
	}

	return exprRecord, nil
}
