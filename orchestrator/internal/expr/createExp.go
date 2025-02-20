package expr

import (
	"fmt"

	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/pkg/db"
)

func CreateExp(exprStore *db.ExpressionStore, taskStore *db.TaskStore, expression string, id string, cfg *config.Config) (db.Expression, error) {
	exprRecord := db.Expression{
		ID:         id,
		RootTaskID: "",
		Status:     db.StatusPending,
		Result:     0,
	}
	exprStore.Create(exprRecord)

	tokens := tokenize(expression)
	parser := NewParser(tokens)
	ast := parser.parseExpression(0)

	rootID, tasks := GenerateTasks(ast, cfg)
	fmt.Printf("Корневой таск ID: %s\n", rootID)

	exprRecord.RootTaskID = rootID
	exprStore.Update(exprRecord)

	for _, task := range tasks {
		if err := taskStore.Create(task); err != nil {
			return exprRecord, err
		}
	}

	return exprRecord, nil
}
