package user

import (
	"app/finance/internal/account"
	"app/finance/internal/card"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex"`
	Password string
	Name     string
	Accounts []account.Account `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"accounts"`
	Cards    []card.Card       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"cards"`
}
