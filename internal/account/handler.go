package account

import (
	"fmt"
	"net/http"
)

type AccountHandler struct{}

type AccountHandlerDeps struct{}

func NewAuthHandler(router *http.ServeMux, deps AccountHandlerDeps) {
	handler := &AccountHandler{}
	router.HandleFunc("POST /account", handler.create())
	router.HandleFunc("GET /account/{id}", handler.read())
	router.HandleFunc("PATCH /account/{id}", handler.update())
	router.HandleFunc("DELETE /account/{id}", handler.delete())
}

func (handler *AccountHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Create account handler")
	}
}

func (handler *AccountHandler) read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Read account handler")
	}
}

func (handler *AccountHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Update account handler")
	}
}

func (handler *AccountHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Delete account handler")
		id := r.PathValue("id")
		fmt.Println("ID: ", id)
	}
}
