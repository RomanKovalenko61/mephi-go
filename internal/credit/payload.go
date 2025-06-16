package credit

import "app/finance/internal/payment"

type CreditCreateRequest struct {
	AccountID uint    `json:"accountId"`
	Amount    float64 `json:"amount"`
	Duration  uint    `json:"duration"` // in months
}

type CreditResponse struct {
	Amount   float64                   `json:"amount"`
	Duration uint                      `json:"duration"`
	Rate     float64                   `json:"rate"`
	Payments []payment.PaymentResponse `json:"payments"`
}
