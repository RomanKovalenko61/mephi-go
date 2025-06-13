package account

import (
	"app/finance/configs"
	"app/finance/internal/card"
	"app/finance/pkg/middleware"
	"app/finance/pkg/request"
	"app/finance/pkg/resp"
	"fmt"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	AccountRepository *AccountRepository
	CardService       *card.CardService
}

type AccountHandlerDeps struct {
	AccountRepository *AccountRepository
	Config            *configs.Config
	CardService       *card.CardService
}

func NewAuthHandler(router *http.ServeMux, deps AccountHandlerDeps) {
	handler := &AccountHandler{
		AccountRepository: deps.AccountRepository,
		CardService:       deps.CardService,
	}
	router.Handle("POST /account", middleware.ISAuthed(handler.create(), deps.Config))
	router.Handle("GET /account", middleware.ISAuthed(handler.getAll(), deps.Config))
	router.Handle("GET /account/{id}", middleware.ISAuthed(handler.read(), deps.Config))
	router.Handle("PATCH /account/{id}", middleware.ISAuthed(handler.update(), deps.Config))
	router.Handle("DELETE /account/{id}", middleware.ISAuthed(handler.delete(), deps.Config))

	router.Handle("POST /cards", middleware.ISAuthed(handler.createCard(), deps.Config))
	router.Handle("GET /cards", middleware.ISAuthed(handler.getAllCards(), deps.Config))
	router.Handle("GET /cards/{id}", middleware.ISAuthed(handler.readCard(), deps.Config))
	router.Handle("DELETE /cards/{id}", middleware.ISAuthed(handler.deleteCard(), deps.Config))
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
		acc, err := handler.getAccountWithUserCheck(w, r)
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
		existedAcc, err := handler.getAccountWithUserCheck(w, r)
		if err != nil {
			return
		}
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
		existedAcc, err := handler.getAccountWithUserCheck(w, r)
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

func (handler *AccountHandler) createCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[card.CardCreateRequest](r)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		userID, err := getIDFromContext(w, r)
		if err != nil {
			return
		}
		acc, err := handler.getAccountById(body.AccountID, userID, w)
		if err != nil {
			return
		}
		newCard, err := handler.CardService.AddCardToAccount(acc.ID, userID)
		if err != nil {
			return
		}
		resp.ResponseJson(w, newCard, http.StatusCreated)
	}
}

func (handler *AccountHandler) readCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathID, err := getPathVariable(w, r)
		if err != nil {
			return
		}
		userID, err := getIDFromContext(w, r)
		if err != nil {
			return
		}
		cardById, err := handler.CardService.GetCardById(userID, pathID)
		if err != nil {
			msg := fmt.Sprintf("Ошибка получения карты с ID %d: %v", pathID, err)
			resp.ResponseJson(w, msg, http.StatusNotFound)
			return
		}
		resp.ResponseJson(w, cardById, http.StatusOK)
	}
}

func (handler *AccountHandler) deleteCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathID, err := getPathVariable(w, r)
		if err != nil {
			return
		}
		userID, err := getIDFromContext(w, r)
		if err != nil {
			return
		}
		deletedCard, err := handler.CardService.DeleteCardById(userID, pathID)
		if err != nil {
			msg := fmt.Sprintf("Карта с ID %d: не доступна %v", pathID, err)
			resp.ResponseJson(w, msg, http.StatusNotFound)
			return
		}
		resp.ResponseJson(w, deletedCard, http.StatusOK)
	}
}

func (handler *AccountHandler) getAllCards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromContext(w, r)
		if err != nil {
			return
		}
		cards, err := handler.CardService.GetAllCards(id)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.ResponseJson(w, cards, http.StatusOK)
	}
}

func getIDFromContext(w http.ResponseWriter, r *http.Request) (uint, error) {
	userID, ok := r.Context().Value(middleware.ContextIDKey).(uint)
	if !ok {
		msg := "не удалось получить ID пользователя из контекста"
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

func (handler *AccountHandler) getAccountWithUserCheck(w http.ResponseWriter, r *http.Request) (*Account, error) {
	id, err := getIDFromContext(w, r)
	if err != nil {
		return nil, err
	}
	pathID, err := getPathVariable(w, r)
	if err != nil {
		return nil, err
	}
	acc, err := handler.getAccountById(pathID, id, w)
	if err != nil {
		msg := fmt.Sprintf("Ошибка получения аккаунта с ID %d: %v", pathID, err)
		resp.ResponseJson(w, msg, http.StatusNotFound)
		return nil, err
	}
	return acc, nil
}

func (handler *AccountHandler) getAccountById(pathID, id uint, w http.ResponseWriter) (*Account, error) {
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
