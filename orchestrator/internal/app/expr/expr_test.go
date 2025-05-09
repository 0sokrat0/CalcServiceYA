package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/entity"
)

func TestTokenize(t *testing.T) {
	input := "2+2*2"
	tokens := tokenize(input)

	assert.Equal(t, TokenNumber, tokens[0].Type)
	assert.Equal(t, "2", tokens[0].Literal)
	assert.Equal(t, TokenPlus, tokens[1].Type)
	assert.Equal(t, "+", tokens[1].Literal)
	assert.Equal(t, TokenNumber, tokens[2].Type)
	assert.Equal(t, "2", tokens[2].Literal)
	assert.Equal(t, TokenMul, tokens[3].Type)
	assert.Equal(t, "*", tokens[3].Literal)
	assert.Equal(t, TokenNumber, tokens[4].Type)
	assert.Equal(t, "2", tokens[4].Literal)
	assert.Equal(t, TokenEOF, tokens[5].Type)
}

func TestParser_ParseExpression(t *testing.T) {
	input := "2+2*2"
	tokens := tokenize(input)
	parser := NewParser(tokens)
	ast := parser.parseExpression(0)

	be, ok := ast.(BinaryExpr)
	assert.True(t, ok, "ast должен быть BinaryExpr")

	left, ok := be.Left.(Number)
	assert.True(t, ok, "левый операнд должен быть Number")
	assert.Equal(t, 2.0, left.Value)

	assert.Equal(t, TokenPlus, be.Operator)

	rightExpr, ok := be.Right.(BinaryExpr)
	assert.True(t, ok, "правый операнд должен быть BinaryExpr")

	rightLeft, ok := rightExpr.Left.(Number)
	assert.True(t, ok, "правый.левый должен быть Number")
	assert.Equal(t, 2.0, rightLeft.Value)

	assert.Equal(t, TokenMul, rightExpr.Operator)

	rightRight, ok := rightExpr.Right.(Number)
	assert.True(t, ok, "правый.правый должен быть Number")
	assert.Equal(t, 2.0, rightRight.Value)
}

func TestGenerateTasks(t *testing.T) {
	ast := BinaryExpr{
		Left: BinaryExpr{
			Left:     Number{Value: 2},
			Operator: TokenPlus,
			Right:    Number{Value: 2},
		},
		Operator: TokenMul,
		Right:    Number{Value: 2},
	}

	cfg := &config.Config{
		App: config.AppConfig{
			TIME_ADDITION_MS:       100,
			TIME_SUBTRACTION_MS:    150,
			TIME_MULTIPLICATION_MS: 200,
			TIME_DIVISION_MS:       250,
		},
	}

	rootID, tasks := GenerateTasks(ast, cfg)
	assert.Equal(t, 2, len(tasks))

	var addTask, mulTask entity.Task
	for _, task := range tasks {
		if task.Operation == "+" {
			addTask = task
		} else if task.Operation == "*" {
			mulTask = task
		}
	}
	assert.NotEmpty(t, addTask.ID)
	assert.Equal(t, "+", addTask.Operation)
	assert.Equal(t, 100, int(addTask.OperationTime))
	assert.NotEmpty(t, mulTask.ID)
	assert.Equal(t, "*", mulTask.Operation)
	assert.Equal(t, 200, int(mulTask.OperationTime))
	assert.Equal(t, mulTask.ID, rootID)
}

func TestOperatorToString(t *testing.T) {
	assert.Equal(t, "+", operatorToString(TokenPlus))
	assert.Equal(t, "-", operatorToString(TokenMinus))
	assert.Equal(t, "*", operatorToString(TokenMul))
	assert.Equal(t, "/", operatorToString(TokenDiv))
	assert.Equal(t, "unknown", operatorToString(999)) // недопустимый тип токена
}

func TestGetOperationTime(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			TIME_ADDITION_MS:       100,
			TIME_SUBTRACTION_MS:    150,
			TIME_MULTIPLICATION_MS: 200,
			TIME_DIVISION_MS:       250,
		},
	}

	assert.Equal(t, int64(100), getOperationTime("+", cfg))
	assert.Equal(t, int64(150), getOperationTime("-", cfg))
	assert.Equal(t, int64(200), getOperationTime("*", cfg))
	assert.Equal(t, int64(250), getOperationTime("/", cfg))
	assert.Equal(t, int64(1000), getOperationTime("unknown", cfg))
}
