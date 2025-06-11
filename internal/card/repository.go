package card

import "app/finance/pkg/db"

type CardRepository struct {
	Database *db.Db
}

func NewAccountRepository(database *db.Db) *CardRepository {
	return &CardRepository{
		Database: database,
	}
}

func (repo *CardRepository) AddCardToAccount(accountId uint, numberEnc, expireEnc, hmac string) (*Card, error) {
	newCard := &Card{
		AccountId: accountId,
		NumberEnc: numberEnc,
		ExpireEnc: expireEnc,
		HMAC:      hmac,
	}
	if err := repo.Database.DB.Table("cards").Create(newCard).Error; err != nil {
		return nil, err
	}
	return newCard, nil
}

func (repo *CardRepository) GetCardById(cardId uint) (*Card, error) {
	var card Card
	if err := repo.Database.DB.Table("cards").First(&card, "id = ?", cardId).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (repo *CardRepository) GetCardsByAccountId(accountId uint) ([]Card, error) {
	var cards []Card
	if err := repo.Database.DB.Table("cards").Where("account_id = ?", accountId).Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}
