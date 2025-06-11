package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

var pgp = crypto.PGP()

func EncryptPGP(data string, password []byte) (string, error) {
	encHandle, err := pgp.Encryption().Password(password).New()
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

func DecryptPGP(enc string, password []byte) (string, error) {
	decHandle, err := pgp.Decryption().Password(password).New()
	if err != nil {
		return "", fmt.Errorf("ошибка создания PGP обработчика: %v", err)
	}
	decrypted, err := decHandle.Decrypt([]byte(enc), crypto.Armor)
	if err != nil {
		return "", fmt.Errorf("ошибка расшифровки PGP сообщения: %v", err)
	}
	return string(decrypted.Bytes()), nil
}

func GenerateCardHMAC(cardNumber, expire, cvv string, password []byte) string {
	data := cardNumber + "|" + expire + "|" + cvv
	mac := hmac.New(sha256.New, password)
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
