package account

import (
	"app/finance/internal/card"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Balance float64   `json:"balance"`
	Owner   string    `json:"owner"`
	Card    card.Card `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"card,omitempty"`
}

func NewAccount(owner string) *Account {
	return &Account{
		Owner: owner,
	}
}
