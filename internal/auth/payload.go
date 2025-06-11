package auth

type LoginRequest struct {
	ID       uint   `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	TOKEN string `json:"token"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	TOKEN string `json:"token"`
}
