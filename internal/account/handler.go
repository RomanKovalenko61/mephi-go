package account

import (
	"app/finance/configs"
	"app/finance/pkg/middleware"
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
	Config            *configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AccountHandlerDeps) {
	handler := &AccountHandler{
		AccountRepository: deps.AccountRepository,
	}
	router.Handle("POST /account", middleware.ISAuthed(handler.create(), deps.Config))
	router.Handle("GET /account", middleware.ISAuthed(handler.getAll(), deps.Config))
	router.Handle("GET /account/{id}", middleware.ISAuthed(handler.read(), deps.Config))
	router.Handle("PATCH /account/{id}", middleware.ISAuthed(handler.update(), deps.Config))
	router.Handle("DELETE /account/{id}", middleware.ISAuthed(handler.delete(), deps.Config))
}

func (handler *AccountHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromContext(w, r)
		if err != nil {
			return
		}
		account := NewAccount(id)
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
		acc, err := getAccountWithUserCheck(w, r, handler)
		if err != nil {
			return
		}
		resp.ResponseJson(w, acc, http.StatusOK)
	}
}

func (handler *AccountHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[AccountUpdateRequest](r)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.Balance < 0 {
			msg := "Баланс не может быть отрицательным"
			resp.ResponseJson(w, msg, http.StatusBadRequest)
			return
		}
		existedAcc, err := getAccountWithUserCheck(w, r, handler)
		existedAcc.Balance = body.Balance
		updateAcc, err := handler.AccountRepository.Update(existedAcc)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp.ResponseJson(w, updateAcc, http.StatusOK)
	}
}

func (handler *AccountHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		existedAcc, err := getAccountWithUserCheck(w, r, handler)
		if existedAcc.Balance != 0 {
			msg := fmt.Sprintf("Невозможно удалить аккаунт с id: %v, так как баланс не равен нулю", existedAcc.ID)
			resp.ResponseJson(w, msg, http.StatusBadRequest)
			return
		}
		err = handler.AccountRepository.Delete(existedAcc.ID)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.ResponseJson(w, "Success deleted", http.StatusOK)
	}
}

func (handler *AccountHandler) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Get all accounts handler")
		id, err := getIDFromContext(w, r)
		if err != nil {
			return
		}
		accounts, err := handler.AccountRepository.GetAll(id)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.ResponseJson(w, accounts, http.StatusOK)
	}
}

func getIDFromContext(w http.ResponseWriter, r *http.Request) (uint, error) {
	userID, ok := r.Context().Value(middleware.ContextIDKey).(uint)
	if ok {
		fmt.Println("Get ID from ctx: ", userID)
	} else {
		msg := "Не удалось получить ID пользователя из контекста"
		resp.ResponseJson(w, msg, http.StatusInternalServerError)
		return 0, fmt.Errorf(msg)
	}
	return userID, nil
}

func getPathVariable(w http.ResponseWriter, r *http.Request) (uint, error) {
	idString := r.PathValue("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		msg := fmt.Sprintf("Ошибка формата ID: %v передано не число", idString)
		resp.ResponseJson(w, msg, http.StatusBadRequest)
		return 0, err
	}
	id := uint(idInt)
	return id, nil
}

func getAccountWithUserCheck(w http.ResponseWriter, r *http.Request, handler *AccountHandler) (*Account, error) {
	id, err := getIDFromContext(w, r)
	if err != nil {
		return nil, err
	}
	pathID, err := getPathVariable(w, r)
	if err != nil {
		return nil, err
	}
	acc, err := handler.AccountRepository.GetById(pathID)
	if err != nil {
		msg := fmt.Sprintf("Аккаунт с id: %v не найден", pathID)
		resp.ResponseJson(w, msg, http.StatusNotFound)
		return nil, err
	}
	if acc.UserID != id {
		msg := fmt.Sprintf("У вас нет прав для доступа к аккаунту с id: %v", pathID)
		resp.ResponseJson(w, msg, http.StatusForbidden)
		return nil, fmt.Errorf(msg)
	}
	return acc, nil
}
