package credit

import (
	"app/finance/internal/account"
	"app/finance/internal/payment"
	"app/finance/pkg/centralbank"
	"app/finance/pkg/db"
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
)

type CreditRepository struct {
	Database *db.Db
}

func NewCreditRepository(database *db.Db) *CreditRepository {
	return &CreditRepository{
		Database: database,
	}
}

func (repo *CreditRepository) create(accountId, userID uint, amount float64, duration uint) (*Credit, error) {
	rate, err := centralbank.GetCentralBankRate()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить данные от ЦБ РФ: %v", err)
	}
	monthlyInterest := rate / 100 / 12
	n := float64(duration)
	r := monthlyInterest
	pmt := amount * ((r * math.Pow(1+r, n)) / (math.Pow(1+r, n) - 1))

	payments := make([]payment.Payment, duration)
	currentDate := time.Now()
	for i := range payments {
		currentDate = currentDate.AddDate(0, 1, 0)
		payments[i].Date = currentDate
		payments[i].Amount = pmt
		payments[i].Redeem = false
	}

	newCredit := &Credit{
		AccountID: accountId,
		UserID:    userID,
		Amount:    amount,
		Rate:      rate,
		Duration:  duration,
	}

	var acc account.Account
	result := repo.Database.DB.Table("accounts").First(&acc, accountId)
	if result.Error != nil {
		return nil, fmt.Errorf("account not found: %v", accountId)
	}
	if acc.UserID != userID {
		return nil, fmt.Errorf("account %d does not belong to user %d", accountId, userID)
	}
	result = repo.Database.DB.Create(&newCredit)
	if result.Error != nil {
		return nil, result.Error
	}

	err = repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		for i := range payments {
			payments[i].ID = 0
			payments[i].CreditID = newCredit.ID
			if err := tx.Table("payments").Create(&payments[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		repo.Database.DB.Delete(newCredit)
		return nil, err
	}
	newCredit.Payments = payments
	return newCredit, nil
}
