package credit

import (
	"app/finance/internal/account"
	"app/finance/internal/payment"
	"app/finance/internal/user"
	"gorm.io/gorm"
)

type Credit struct {
	gorm.Model
	UserID    uint            `json:"userId" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User      user.User       `gorm:"foreignKey:UserID;references:ID"`
	AccountID uint            `json:"accountId" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Account   account.Account `gorm:"foreignKey:AccountID;references:ID"`
	Amount    float64
	Rate      float64
	Duration  uint
	Payments  []payment.Payment `gorm:"foreignKey:CreditID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
