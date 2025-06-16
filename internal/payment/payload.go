package payment

type PaymentResponse struct {
	Order  uint   `json:"order"`
	Date   string `json:"date"`
	Amount string `json:"amount"`
}
