package cardutil

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// https://habr.com/ru/articles/745302/ как проверить номер карты по алгоритму Луна
// https://go.dev/play/p/f-v7_OouV68 тестировал генерацию и проверку номера карты

func GenerateCardNumber() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	base := make([]int, 15)
	for i := range base {
		base[i] = rnd.Intn(10)
	}
	sum := 0
	for i := 0; i < 15; i++ {
		if i%2 == 0 {
			sum += base[i]
		} else {
			doubled := base[i] * 2
			if doubled > 9 {
				doubled -= 9
			}
			sum += doubled
		}
	}
	checkDigit := (10 - (sum % 10)) % 10
	cardNumber := ""
	for _, digit := range base {
		cardNumber += strconv.Itoa(digit)
	}
	cardNumber += strconv.Itoa(checkDigit)
	return cardNumber
}

func CheckCardNumber(cardNumber string) bool {
	cardDigits := strings.Split(strings.ReplaceAll(cardNumber, " ", ""), "")
	if len(cardDigits) != 16 {
		return false
	}

	var sum int
	for i := 0; i < 16; i++ {
		digit, _ := strconv.Atoi(cardDigits[i])
		if i%2 != 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
