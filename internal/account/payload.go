package account

type AccountCreateRequest struct {
	Owner string `json:"owner" validate:"required"`
}

type AccountUpdateRequest struct {
	Balance float64 `json:"balance" validate:"required"`
}
