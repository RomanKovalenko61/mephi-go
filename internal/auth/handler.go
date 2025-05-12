package auth

import (
	"fmt"
	"net/http"
)

type AuthHandler struct{}

func NewAuthHandler(router *http.ServeMux) {
	handler := &AuthHandler{}
	router.HandleFunc("POST /auth/login", handler.login())
	router.HandleFunc("POST /auth/register", handler.register())
}

func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login handler")
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("registrer handler")
	}
}
