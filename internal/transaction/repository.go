package transaction

import (
	"app/finance/internal/account"
	"app/finance/pkg/db"
	"fmt"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	Database *db.Db
}

func NewTransactionRepository(database *db.Db) *TransactionRepository {
	return &TransactionRepository{
		Database: database,
	}
}

func (repo *TransactionRepository) Transfer(fromID, toID uint, amount float64) error {
	return repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		var fromAccount, toAccount account.Account
		if err := tx.First(&fromAccount, fromID).Error; err != nil {
			return fmt.Errorf("account not found: %v", fromID)
		}
		if err := tx.First(&toAccount, toID).Error; err != nil {
			return fmt.Errorf("account not found: %v", toID)
		}

		if fromAccount.Balance < amount {
			return fmt.Errorf("Недостаточно средств на счете %d для перевода %f", fromID, amount)
		}

		fromAccount.Balance -= amount
		toAccount.Balance += amount

		if err := tx.Save(&fromAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&toAccount).Error; err != nil {
			return err
		}

		transaction := Transaction{
			FromID: fromAccount.ID,
			ToID:   toAccount.ID,
			Amount: amount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		return nil
	})
}

func (repo *TransactionRepository) FindById(id uint) (*Transaction, error) {
	var transaction Transaction
	result := repo.Database.DB.Table("users").First(&transaction, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

func (repo *TransactionRepository) GetAllByUserId(userID uint) ([]Transaction, error) {
	var transactions []Transaction
	result := repo.Database.DB.Table("transactions").Where("from = ? OR to = ?", userID, userID).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (repo *TransactionRepository) GetAll() ([]Transaction, error) {
	var transactions []Transaction
	result := repo.Database.DB.Table("transactions").Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (repo *TransactionRepository) Delete(id uint) error {
	result := repo.Database.DB.Table("transactions").Delete(&Transaction{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
