package transaction

import (
	"app/finance/internal/account"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	FromID uint            `json:"from" gorm:"not null"`
	From   account.Account `gorm:"foreignKey:FromID;references:ID"`
	ToID   uint            `json:"to" gorm:"not null"`
	To     account.Account `gorm:"foreignKey:ToID;references:ID"`
	Amount float64         `json:"amount" gorm:"not null"`
}
