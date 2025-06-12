package user

import (
	"app/finance/internal/account"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"index"`
	Password string
	Name     string
	Accounts []account.Account `gorm:"foreignKey:UserID"`
}
