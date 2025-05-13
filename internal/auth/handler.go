package auth

import (
	"app/finance/configs"
	"app/finance/pkg/request"
	"app/finance/pkg/resp"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	*configs.Config
	*AuthService
}

// TODO: excess deps
type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.login())
	router.HandleFunc("POST /auth/register", handler.register())
}

func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login handler")
		body, err := request.HandleBody[LoginRequest](r)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Payload: ", body)
		data := LoginResponse{
			TOKEN: "123",
		}
		resp.ResponseJson(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("registrer handler")
		body, err := request.HandleBody[RegisterRequest](r)
		if err != nil {
			resp.ResponseJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Payload: ", body)
		handler.AuthService.Register(body.Email, body.Password, body.Username)
		data := RegisterResponse{
			TOKEN: "123",
		}
		resp.ResponseJson(w, data, http.StatusOK)
	}
}
