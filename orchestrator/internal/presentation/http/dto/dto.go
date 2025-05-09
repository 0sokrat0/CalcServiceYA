package dto

import "github.com/golang-jwt/jwt/v5"

type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type CalculateRequest struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

type TasksResponse struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
}
type ExpressionResponse struct {
	ID         string  `json:"id"`
	RootTaskID string  `json:"root_task_id"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type,omitempty"`
}
