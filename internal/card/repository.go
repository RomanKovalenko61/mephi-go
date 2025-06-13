package card

import "app/finance/pkg/db"

type CardRepository struct {
	Database *db.Db
}

func NewCardRepository(database *db.Db) *CardRepository {
	return &CardRepository{
		Database: database,
	}
}

func (repo *CardRepository) AddCardToAccount(accountId uint, userID uint, numberEnc, expireEnc, cvv, hmac string) (*Card, error) {
	newCard := &Card{
		AccountID: accountId,
		UserID:    userID,
		NumberEnc: numberEnc,
		ExpireEnc: expireEnc,
		CVV:       cvv,
		HMAC:      hmac,
	}
	result := repo.Database.DB.Table("cards").Create(newCard)
	if result.Error != nil {
		return nil, result.Error
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

func (repo *CardRepository) DeleteCardById(cardID uint) error {
	result := repo.Database.DB.Table("cards").Delete(&Card{}, cardID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *CardRepository) GetAll(userID uint) ([]Card, error) {
	var cards []Card
	result := repo.Database.DB.Table("cards").Where("user_id = ?", userID).Find(&cards)
	if result.Error != nil {
		return nil, result.Error
	}
	return cards, nil
}
