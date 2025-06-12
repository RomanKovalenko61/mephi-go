package account

import (
	"app/finance/pkg/db"

	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	Database *db.Db
}

func NewAccountRepository(database *db.Db) *AccountRepository {
	return &AccountRepository{
		Database: database,
	}
}

func (repo *AccountRepository) Create(acc *Account) (*Account, error) {
	result := repo.Database.DB.Table("accounts").Create(acc)
	if result.Error != nil {
		return nil, result.Error
	}
	return acc, nil
}

func (repo *AccountRepository) GetById(id uint) (*Account, error) {
	var acc Account
	result := repo.Database.DB.Table("accounts").First(&acc, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &acc, nil
}

func (repo *AccountRepository) Update(acc *Account) (*Account, error) {
	result := repo.Database.DB.Model(&acc).Table("accounts").
		Clauses(clause.Returning{}).
		Omit(clause.Associations).
		Updates(map[string]interface{}{"balance": acc.Balance})

	if result.Error != nil {
		return nil, result.Error
	}
	return acc, nil
}

func (repo *AccountRepository) Delete(id uint) error {
	result := repo.Database.DB.Table("accounts").Delete(&Account{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *AccountRepository) GetAll(id uint) ([]Account, error) {
	var accounts []Account
	result := repo.Database.DB.Table("accounts").Where("user_id = ?", id).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}
