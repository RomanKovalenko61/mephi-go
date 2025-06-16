package transaction

import (
	"app/finance/configs"
	"app/finance/pkg/middleware"
	"app/finance/pkg/request"
	"app/finance/pkg/resp"
	"net/http"
)

type TransactionHandler struct {
	TransactionRepository *TransactionRepository
}

type TransactionHandlerDeps struct {
	TransactionRepository *TransactionRepository
	Config                *configs.Config
}

func NewTransactionHandler(router *http.ServeMux, deps TransactionHandlerDeps) {
	handler := &TransactionHandler{
		TransactionRepository: deps.TransactionRepository,
	}
	router.Handle("POST /transfer", middleware.ISAuthed(handler.create(), deps.Config))
}

func (handler *TransactionHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[TransactionRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.TransactionRepository.Transfer(body.FromID, body.ToID, body.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp.ResponseJson(w, "success", http.StatusCreated)
	}
}
