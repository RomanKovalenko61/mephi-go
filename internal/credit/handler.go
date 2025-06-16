package credit

import (
	"app/finance/configs"
	"app/finance/internal/payment"
	"app/finance/pkg/middleware"
	"app/finance/pkg/request"
	"app/finance/pkg/resp"
	"fmt"
	"net/http"
)

type CreditHandler struct {
	CreditRepository *CreditRepository
}

type CreditHandlerDeps struct {
	CreditRepository *CreditRepository
	Config           *configs.Config
}

func NewCreditHandler(router *http.ServeMux, deps CreditHandlerDeps) {
	handler := &CreditHandler{
		CreditRepository: deps.CreditRepository,
	}
	router.Handle("POST /credit", middleware.ISAuthed(handler.create(), deps.Config))
}

func (handler *CreditHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[CreditCreateRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userID, err := getIDFromContext(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		credit, err := handler.CreditRepository.create(body.AccountID, userID, body.Amount, body.Duration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		payments := make([]payment.PaymentResponse, len(credit.Payments))
		for i := range payments {
			payments[i] = payment.PaymentResponse{
				Order:  uint(i + 1),
				Date:   credit.Payments[i].Date.Format("02-01-2006"),
				Amount: fmt.Sprintf("%.2f", credit.Payments[i].Amount),
			}
		}
		response := CreditResponse{
			Amount:   credit.Amount,
			Duration: credit.Duration,
			Rate:     credit.Rate,
			Payments: payments,
		}
		resp.ResponseJson(w, response, http.StatusCreated)
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
