package card

import (
	"app/finance/pkg/cardutil"
	"app/finance/pkg/crypto"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type CardService struct {
	CardRepository *CardRepository
	CryptoHelper   *crypto.CryptoHelper
}

type CardServiceDeps struct {
	CardRepository *CardRepository
	CryptoHelper   *crypto.CryptoHelper
}

func NewCardService(deps CardServiceDeps) *CardService {
	return &CardService{
		CardRepository: deps.CardRepository,
		CryptoHelper:   deps.CryptoHelper,
	}
}

func (service *CardService) AddCardToAccount(accountId, userID uint) (*Card, error) {
	cardNumber := cardutil.GenerateCardNumber()
	expireDate := time.Now().AddDate(3, 0, 0).String()
	numberEnc, err := service.CryptoHelper.EncryptPGP(cardNumber)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования номера карты: %w", err)
	}
	expireEnc, err := service.CryptoHelper.EncryptPGP(expireDate)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования даты истечения карты: %w", err)
	}
	cvvRow := cardutil.GenerateCVV()
	cvv, err := service.CryptoHelper.EncryptPGP(cvvRow)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования CVV: %w", err)
	}
	hmac := service.CryptoHelper.GenerateCardHMAC(cardNumber, expireDate, cvvRow)
	savedCard, err := service.CardRepository.AddCardToAccount(accountId, userID, numberEnc, expireEnc, cvv, hmac)
	if err != nil {
		return nil, err
	}
	return &Card{
		Model:     gorm.Model{ID: savedCard.ID},
		AccountID: accountId, UserID: userID,
		NumberEnc: cardNumber,
		ExpireEnc: expireDate,
		CVV:       cvvRow,
		HMAC:      hmac}, nil
}

func (service *CardService) GetCardById(userID uint, cardID uint) (*Card, error) {
	existedCard, err := service.checkAccessAndReturnCard(userID, cardID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения карты по ID: %w", err)
	}
	numberEnc, err := service.CryptoHelper.DecryptPGP(existedCard.NumberEnc)
	if err != nil {
		return nil, fmt.Errorf("ошибка расшифровки номера карты: %w", err)
	}
	expireEnc, err := service.CryptoHelper.DecryptPGP(existedCard.ExpireEnc)
	if err != nil {
		return nil, fmt.Errorf("ошибка расшифровки даты истечения карты: %w", err)
	}
	cvv, err := service.CryptoHelper.DecryptPGP(existedCard.CVV)
	if err != nil {
		return nil, fmt.Errorf("ошибка расшифровки CVV: %w", err)
	}
	encryptCard := &Card{
		Model:     gorm.Model{ID: existedCard.ID},
		AccountID: existedCard.AccountID,
		UserID:    existedCard.UserID,
		NumberEnc: numberEnc,
		ExpireEnc: expireEnc,
		HMAC:      existedCard.HMAC,
		CVV:       cvv,
	}
	return encryptCard, nil
}

func (service *CardService) DeleteCardById(userID uint, cardID uint) (*Card, error) {
	existedCard, err := service.checkAccessAndReturnCard(userID, cardID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения карты по ID: %w", err)
	}
	err = service.CardRepository.DeleteCardById(existedCard.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка удаления карты: %w", err)
	}
	return existedCard, nil
}

func (service *CardService) GetAllCards(userID uint) ([]Card, error) {
	cards, err := service.CardRepository.GetAll(userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех карт: %w", err)
	}
	return cards, nil
}

func (service *CardService) checkAccessAndReturnCard(userID uint, cardID uint) (*Card, error) {
	card, err := service.CardRepository.GetCardById(cardID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения карты по ID: %w", err)
	}
	if card.UserID != userID {
		return nil, fmt.Errorf("карта не принадлежит пользователю с ID %d", userID)
	}
	return card, nil
}
