package account

import (
	"app/finance/pkg/request"
	"app/finance/pkg/resp"
	"fmt"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	AccountRepository *AccountRepository
}

type AccountHandlerDeps struct {
	AccountRepository *AccountRepository
}

func NewAuthHandler(router *http.ServeMux, deps AccountHandlerDeps) {
	handler := &AccountHandler{
		AccountRepository: deps.AccountRepository,
	}
	router.HandleFunc("POST /account", handler.create())
	router.HandleFunc("GET /account/{id}", handler.read())
	router.HandleFunc("PATCH /account/{id}", handler.update())
	router.HandleFunc("DELETE /account/{id}", handler.delete())
}

func (handler *AccountHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Create account handler")
		body, err := request.HandleBody[AccountCreateRequest](r)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		account := NewAccount(body.Owner)
		createdAcc, err := handler.AccountRepository.Create(account)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp.ResponseJson(w, createdAcc, http.StatusCreated)
	}
}

func (handler *AccountHandler) read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Read account handler")
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			msg := fmt.Sprintf("id: %v isn't number", idString)
			resp.ResponseJson(w, msg, http.StatusBadRequest)
			return
		}
		acc, err := handler.AccountRepository.GetById(id)
		if err != nil {
			msg := fmt.Sprintf("Not Found Account with id: %v", id)
			resp.ResponseJson(w, msg, http.StatusNotFound)
			return
		}
		resp.ResponseJson(w, acc, http.StatusOK)
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
