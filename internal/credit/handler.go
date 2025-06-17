package credit

import (
	"app/finance/configs"
	"app/finance/internal/payment"
	"app/finance/pkg/middleware"
	"app/finance/pkg/request"
	"app/finance/pkg/resp"
	"fmt"
	"net/http"
	"strconv"
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
	router.Handle("GET /credit/{id}/schedule", middleware.ISAuthed(handler.get(), deps.Config))
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
		response := paymentsToResponse(credit)
		resp.ResponseJson(w, response, http.StatusCreated)
	}
}

func (handler *CreditHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		credit, err := handler.getCreditWithUserCheck(w, r)
		if err != nil {
			return
		}
		response := paymentsToResponse(credit)
		resp.ResponseJson(w, response, http.StatusOK)
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

func (handler *CreditHandler) getCreditWithUserCheck(w http.ResponseWriter, r *http.Request) (*Credit, error) {
	id, err := getIDFromContext(w, r)
	if err != nil {
		return nil, err
	}
	pathID, err := getPathVariable(w, r)
	if err != nil {
		return nil, err
	}
	credit, err := handler.getCreditById(pathID, id, w)
	if err != nil {
		msg := fmt.Sprintf("Ошибка получения кредита с ID %d: %v", pathID, err)
		resp.ResponseJson(w, msg, http.StatusNotFound)
		return nil, err
	}
	return credit, nil
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

func (handler *CreditHandler) getCreditById(pathID, id uint, w http.ResponseWriter) (*Credit, error) {
	credit, err := handler.CreditRepository.GetById(pathID)
	if err != nil {
		msg := fmt.Sprintf("Кредит с id: %v не найден", pathID)
		resp.ResponseJson(w, msg, http.StatusNotFound)
		return nil, err
	}
	if credit.UserID != id {
		msg := fmt.Sprintf("У вас нет прав для доступа к кредиту с id: %v", pathID)
		resp.ResponseJson(w, msg, http.StatusForbidden)
		return nil, fmt.Errorf(msg)
	}
	return credit, nil
}

func paymentsToResponse(credit *Credit) *CreditResponse {
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
	return &response
}
