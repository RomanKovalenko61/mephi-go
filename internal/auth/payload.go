package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	TOKEN string `json:"token"`
}
