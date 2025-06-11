package card

import (
	"app/finance/pkg/cardutil"
	"time"
)

type CardService struct {
	CardRepository *CardRepository
}

func NewCardService(cardRepository *CardRepository) *CardService {
	return &CardService{
		CardRepository: cardRepository,
	}
}

func (service *CardService) AddCardToAccount(accountId uint) (*Card, error) {
	cardNumber := cardutil.GenerateCardNumber()
	expireDate := time.Now().AddDate(3, 0, 0).String()
	newCard, err := service.CardRepository.AddCardToAccount(accountId, cardNumber, expireDate, "")
	if err != nil {
		return nil, err
	}
	return newCard, nil
}
