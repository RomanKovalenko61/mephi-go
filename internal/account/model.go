package account

import (
	"app/finance/internal/card"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID  uint        `json:"userId"`
	Balance float64     `json:"balance"`
	Cards   []card.Card `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

func NewAccount(userID uint) *Account {
	return &Account{
		UserID: userID,
	}
}
