package account

type AccountCreateRequest struct {
	Owner string `json:"owner" validate:"required"`
}
