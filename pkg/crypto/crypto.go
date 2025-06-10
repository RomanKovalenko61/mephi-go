package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"io"
)

func EncryptPGP(data string, pubKey []byte) (string, error) {
	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(pubKey))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	w, err := openpgp.Encrypt(&buf, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}
	defer func(w io.WriteCloser) {
		err := w.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия шифратора: %v\n", err)
		}
	}(w)

	_, err = w.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func DecryptPGP(enc string, privKey []byte) (string, error) {
	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(privKey))
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString(enc)
	md, err := openpgp.ReadMessage(buf, entityList, nil, nil)
	if err != nil {
		return "", err
	}
	decryptedData := new(bytes.Buffer)
	_, err = decryptedData.ReadFrom(md.UnverifiedBody)
	if err != nil {
		return "", err
	}
	return decryptedData.String(), nil
}

func ComputeHMAC(data string, secret []byte) (string, error) {
	h := hmac.New(sha256.New, secret)

	_, err := h.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("Ошибка записи данных в HMAC: %v", err)
	}
	result := hex.EncodeToString(h.Sum(nil))
	return result, nil
}
