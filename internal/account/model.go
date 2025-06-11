package account

import (
	"app/finance/internal/card"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID  uint      `json:"userId"`
	Balance float64   `json:"balance"`
	Card    card.Card `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"card,omitempty"`
}

func NewAccount(userID uint) *Account {
	return &Account{
		UserID: userID,
	}
}
