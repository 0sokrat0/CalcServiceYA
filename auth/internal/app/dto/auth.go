package dto

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
