package main

import (
	"app/finance/internal/account"
	"app/finance/internal/card"
	"app/finance/internal/credit"
	"app/finance/internal/payment"
	"app/finance/internal/transaction"
	"app/finance/internal/user"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&user.User{}, &account.Account{}, &card.Card{}, &transaction.Transaction{},
		&credit.Credit{}, &payment.Payment{})
}
