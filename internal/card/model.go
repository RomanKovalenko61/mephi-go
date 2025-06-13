package card

import (
	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	AccountID uint   `json:"accountId" gorm:"foreignKey:AccountID;references:ID"`
	UserID    uint   `json:"userId" gorm:"foreignKey:UserID;references:ID"`
	NumberEnc string `json:"numberEnc"`
	ExpireEnc string `json:"expireEnc"`
	CVV       string `json:"cvv"`
	HMAC      string `json:"hmac"`
}
