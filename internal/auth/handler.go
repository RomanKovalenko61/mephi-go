package auth

import (
	"app/finance/configs"
	"app/finance/pkg/resp"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	*configs.Config
}

// TODO: excess deps
type AuthHandlerDeps struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.login())
	router.HandleFunc("POST /auth/register", handler.register())
}

func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login handler")
		var payload LoginRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			resp.ResponseJson(w, err, http.StatusBadRequest)
			return
		}
		if payload.Email == "" {
			resp.ResponseJson(w, "Email required", http.StatusBadRequest)
			return
		}
		if payload.Password == "" {
			resp.ResponseJson(w, "Password required", http.StatusBadRequest)
			return
		}
		fmt.Println("Payload: ", payload)
		data := LoginResponse{
			TOKEN: "123",
		}
		resp.ResponseJson(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("registrer handler")
	}
}
