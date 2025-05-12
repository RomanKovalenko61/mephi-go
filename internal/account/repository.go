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

func (repo *AccountRepository) GetById(id int) (*Account, error) {
	var acc Account
	result := repo.Database.DB.Table("accounts").First(&acc, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &acc, nil
}

func (repo *AccountRepository) Update(acc *Account) (*Account, error) {
	result := repo.Database.DB.Table("accounts").Clauses(clause.Returning{}).Updates(acc)
	if result.Error != nil {
		return nil, result.Error
	}
	return acc, nil
}
