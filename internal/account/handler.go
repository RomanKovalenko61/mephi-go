package account

import (
	"app/finance/configs"
	"app/finance/pkg/middleware"
	"app/finance/pkg/request"
	"app/finance/pkg/resp"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type AccountHandler struct {
	AccountRepository *AccountRepository
}

type AccountHandlerDeps struct {
	AccountRepository *AccountRepository
	Config            *configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AccountHandlerDeps) {
	handler := &AccountHandler{
		AccountRepository: deps.AccountRepository,
	}
	router.Handle("POST /account", middleware.ISAuthed(handler.create(), deps.Config))
	router.Handle("GET /account/{id}", middleware.ISAuthed(handler.read(), deps.Config))
	router.Handle("PATCH /account/{id}", middleware.ISAuthed(handler.update(), deps.Config))
	router.Handle("DELETE /account/{id}", middleware.ISAuthed(handler.delete(), deps.Config))
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
		if email, ok := r.Context().Value(middleware.ContextEmailKey).(string); ok {
			fmt.Println("Get Email from ctx: ", email)
		}
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			msg := fmt.Sprintf("id: %v isn't number", idString)
			resp.ResponseJson(w, msg, http.StatusBadRequest)
			return
		}
		acc, err := handler.AccountRepository.GetById(uint(id))
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
		body, err := request.HandleBody[AccountUpdateRequest](r)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			msg := fmt.Sprintf("id: %v isn't number", idString)
			resp.ResponseJson(w, msg, http.StatusBadRequest)
			return
		}
		acc, err := handler.AccountRepository.Update(&Account{
			Model: gorm.Model{
				ID: uint(id),
			},
			Balance: body.Balance,
		})
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp.ResponseJson(w, acc, http.StatusOK)
	}
}

func (handler *AccountHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Delete account handler")
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			msg := fmt.Sprintf("id: %v isn't number", idString)
			resp.ResponseJson(w, msg, http.StatusBadRequest)
			return
		}
		_, err = handler.AccountRepository.GetById(uint(id))
		if err != nil {
			resp.ResponseJson(w, "Wrong id", http.StatusNotFound)
			return
		}
		err = handler.AccountRepository.Delete(uint(id))
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.ResponseJson(w, "Success deleted", http.StatusOK)
	}
}
