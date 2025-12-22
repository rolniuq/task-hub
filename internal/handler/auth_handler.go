package handler

import (
	"encoding/json"
	"fmt"
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

func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func writeHTMXError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<div class="alert alert-error shake">%s</div>`, message)
}

func writeHTMXSuccess(w http.ResponseWriter, message string, redirectURL string) {
	w.Header().Set("Content-Type", "text/html")
	if redirectURL != "" {
		w.Header().Set("HX-Redirect", redirectURL)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<div class="alert alert-success">%s</div>`, message)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	isHTMX := isHTMXRequest(r)

	var req app.RegisterRequest
	if isHTMX {
		req.Name = r.FormValue("name")
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		if isHTMX {
			writeHTMXError(w, "Email, password, and name are required")
			return
		}
		writeError(w, http.StatusBadRequest, "email, password, and name are required")
		return
	}

	resp, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		if err == app.ErrUserAlreadyExists {
			if isHTMX {
				writeHTMXError(w, "An account with this email already exists")
				return
			}
			writeError(w, http.StatusConflict, "user already exists")
			return
		}
		if isHTMX {
			writeHTMXError(w, "Failed to create account. Please try again.")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	if isHTMX {
		writeHTMXSuccess(w, "Account created successfully! Redirecting to login...", "/login")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	isHTMX := isHTMXRequest(r)

	var req app.LoginRequest
	if isHTMX {
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
	}

	if req.Email == "" || req.Password == "" {
		if isHTMX {
			writeHTMXError(w, "Email and password are required")
			return
		}
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if err == app.ErrInvalidCredentials {
			if isHTMX {
				writeHTMXError(w, "Invalid email or password")
				return
			}
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if isHTMX {
			writeHTMXError(w, "Login failed. Please try again.")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to login")
		return
	}

	if isHTMX {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    resp.Tokens.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			MaxAge:   900,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.Tokens.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			MaxAge:   604800,
		})
		writeHTMXSuccess(w, "Login successful! Redirecting...", "/dashboard")
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

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
