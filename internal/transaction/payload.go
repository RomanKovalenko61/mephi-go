package transaction

type TransactionRequest struct {
	FromID uint    `json:"fromID" validate:"required"`
	ToID   uint    `json:"toID" validate:"required"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}
