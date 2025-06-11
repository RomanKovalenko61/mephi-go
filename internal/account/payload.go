package account

type AccountCreateRequest struct {
	UserID uint `json:"userID" validate:"required"`
}

type AccountUpdateRequest struct {
	Balance float64 `json:"balance" validate:"required"`
}
