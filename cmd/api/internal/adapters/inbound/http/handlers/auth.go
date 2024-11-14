package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/japablazatww/song-searcher/cmd/api/internal/application/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type CreateAppRequest struct {
	AppName string `json:"app_name"`
}

type GenerateTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (h *AuthHandler) CreateApp(w http.ResponseWriter, r *http.Request) {
	var req CreateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	app, err := h.authService.CreateApp(req.AppName)
	if err != nil {
		http.Error(w, "Error creating app", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(app)
}

func (h *AuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	var req GenerateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(req.ClientID, req.ClientSecret)
	if err != nil {
		fmt.Println("DESDE EL HANDLER ", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
