package calc

import "CalcYA/orchestrator/pkg/db"

func CreateExp(store *db.ExpressionStore, expression string, id string) (db.Expression, error) {
	expr := db.Expression{
		ID:     id,
		Status: db.StatusAccepted,
		Result: 0,
	}
	store.Create(expr)
	return expr, nil
}
