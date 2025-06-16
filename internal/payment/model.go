package payment

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	CreditID uint      `json:"credit_id"`
	Date     time.Time `json:"date"`
	Amount   float64   `json:"amount"`
	Redeem   bool      `json:"redeem"`
}
