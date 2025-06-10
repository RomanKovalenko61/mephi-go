package card

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	AccountId uint `json:"accountId"`
	NumberEnc string
	ExpireEnc string
	HMAC      string
}
