package account

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Balance float64 `json:"balance"`
	Owner   string  `json:"owner"`
}

func NewAccount(owner string) *Account {
	return &Account{
		Owner: owner,
	}
}
