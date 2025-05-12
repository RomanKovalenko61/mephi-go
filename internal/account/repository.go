package account

import "app/finance/pkg/db"

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
