package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret string
}

type JWTData struct {
	UserID    uint
	IssuedAt  int64
	ExpiresAt int64
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data JWTData) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    data.UserID,
		"issuedAt":  jwt.NewNumericDate(time.Now()),
		"expiresAt": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil || !t.Valid {
		return false, nil
	}
	expiresAtParsed, ok := t.Claims.(jwt.MapClaims)["expiresAt"]
	if !ok {
		return false, nil
	}
	expiresAtInt, ok := expiresAtParsed.(float64)
	if !ok {
		return false, nil
	}
	expiresAt := int64(expiresAtInt)
	expirationTime := time.Unix(expiresAt, 0)
	now := time.Now()
	if now.After(expirationTime) {
		return false, nil
	}
	userID := t.Claims.(jwt.MapClaims)["userID"]
	return t.Valid, &JWTData{
		UserID: uint(userID.(float64)),
	}
}
