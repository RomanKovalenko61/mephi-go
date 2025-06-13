package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

type CryptoHelper struct {
	secretKey string
}

func NewCryptoHelper(secretKey string) *CryptoHelper {
	return &CryptoHelper{
		secretKey: secretKey,
	}
}

var pgp = crypto.PGP()

func (helper *CryptoHelper) EncryptPGP(data string) (string, error) {
	encHandle, err := pgp.Encryption().Password([]byte(helper.secretKey)).New()
	if err != nil {
		return "", fmt.Errorf("ошибка создания PGP обработчика: %v", err)
	}
	pgpMessage, err := encHandle.Encrypt([]byte(data))
	if err != nil {
		return "", fmt.Errorf("ошибка создания PGP сообщения: %v", err)
	}
	armored, err := pgpMessage.ArmorBytes()
	if err != nil {
		return "", fmt.Errorf("ошибка кодирования PGP сообщения: %v", err)
	}
	return string(armored), nil
}

func (helper *CryptoHelper) DecryptPGP(enc string) (string, error) {
	decHandle, err := pgp.Decryption().Password([]byte(helper.secretKey)).New()
	if err != nil {
		return "", fmt.Errorf("ошибка создания PGP обработчика: %v", err)
	}
	decrypted, err := decHandle.Decrypt([]byte(enc), crypto.Armor)
	if err != nil {
		return "", fmt.Errorf("ошибка расшифровки PGP сообщения: %v", err)
	}
	return string(decrypted.Bytes()), nil
}

func (helper *CryptoHelper) GenerateCardHMAC(cardNumber, expire, cvv string) string {
	data := cardNumber + "|" + expire + "|" + cvv
	mac := hmac.New(sha256.New, []byte(helper.secretKey))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
