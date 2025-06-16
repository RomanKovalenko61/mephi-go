package main

import (
	"app/finance/configs"
	"app/finance/internal/account"
	"app/finance/internal/auth"
	"app/finance/internal/card"
	"app/finance/internal/credit"
	"app/finance/internal/transaction"
	"app/finance/internal/user"
	"app/finance/pkg/crypto"
	"app/finance/pkg/db"
	"app/finance/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)

	router := http.NewServeMux()

	//Helper
	cryptoHelper := crypto.NewCryptoHelper(conf.Crypto.Secret)

	//Repositories
	accountRepository := account.NewAccountRepository(db)
	userRepository := user.NewUserRepository(db)
	cardRepository := card.NewCardRepository(db)
	transactionRepository := transaction.NewTransactionRepository(db)
	creditRepository := credit.NewCreditRepository(db)

	//Services
	authService := auth.NewAuthService(userRepository)

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	account.NewAuthHandler(router, account.AccountHandlerDeps{
		AccountRepository: accountRepository,
		Config:            conf,
		CardService: card.NewCardService(card.CardServiceDeps{
			CardRepository: cardRepository,
			CryptoHelper:   cryptoHelper,
		}),
	})
	transaction.NewTransactionHandler(router, transaction.TransactionHandlerDeps{
		TransactionRepository: transactionRepository,
		Config:                conf,
	})
	credit.NewCreditHandler(router, credit.CreditHandlerDeps{
		CreditRepository: creditRepository,
		Config:           conf,
	})

	//Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
