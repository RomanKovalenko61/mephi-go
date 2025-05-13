package main

import (
	"app/finance/configs"
	"app/finance/internal/account"
	"app/finance/internal/auth"
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
	accountReposotory := account.NewAccountRepository(db)

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	account.NewAuthHandler(router, account.AccountHandlerDeps{
		AccountRepository: accountReposotory,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: middleware.CORS(middleware.Logging(router)),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
