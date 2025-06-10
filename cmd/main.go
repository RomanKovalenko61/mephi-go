package main

import (
	"app/finance/configs"
	"app/finance/internal/account"
	"app/finance/internal/auth"
	"app/finance/internal/user"
	"app/finance/pkg/db"
	"app/finance/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)

	router := http.NewServeMux()

	//Repositories
	accountRepository := account.NewAccountRepository(db)
	userRepository := user.NewUserRepository(db)

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
