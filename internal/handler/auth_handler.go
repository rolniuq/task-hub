package handler

import (
	"encoding/json"
	"net/http"
	"taskhub/internal/app"
)

type AuthHandler struct {
	authService *app.AuthService
}

func NewAuthHandler(authService *app.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req app.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "email, password, and name are required")
		return
	}

	resp, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		if err == app.ErrUserAlreadyExists {
			writeError(w, http.StatusConflict, "user already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req app.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if err == app.ErrInvalidCredentials {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to login")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req app.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokens, err := h.authService.RefreshToken(r.Context(), &req)
	if err != nil {
		if err == app.ErrTokenExpired {
			writeError(w, http.StatusUnauthorized, "refresh token expired")
			return
		}
		if err == app.ErrInvalidToken {
			writeError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to refresh token")
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}
