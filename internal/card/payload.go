package card

type CardCreateRequest struct {
	AccountID uint `json:"accountID" validate:"required"`
}
