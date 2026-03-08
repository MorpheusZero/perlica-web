package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/morpheuszero/perlica-web/config"
	"github.com/morpheuszero/perlica-web/server/database/repositories"
	"github.com/morpheuszero/perlica-web/server/guards"
	"github.com/morpheuszero/perlica-web/server/services"
	"github.com/morpheuszero/perlica-web/server/util"
)

type AuthController struct {
	authGuard   *guards.AuthGuard
	config      *config.EnvProvider
	authService *services.AuthService
	userService *services.UserService
}

func NewAuthController(authGuard *guards.AuthGuard, config *config.EnvProvider, authService *services.AuthService, userService *services.UserService) *AuthController {
	return &AuthController{
		authGuard:   authGuard,
		config:      config,
		authService: authService,
		userService: userService,
	}
}

func (c *AuthController) MapController() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/login", c.login)
	router.With(c.authGuard.ValidateSession).Get("/session", c.getSession)

	return router
}

func (c *AuthController) login(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusBadRequest)
		return
	}

	// Expecting "Basic base64(username:password)"
	const prefix = "Basic "
	if !util.StartsWith(authHeader, prefix) {
		http.Error(w, "Invalid authorization header format", http.StatusBadRequest)
		return
	}

	username, password, err := util.DecodeBasicAuth(authHeader[len(prefix):])
	if err != nil {
		http.Error(w, "Invalid authorization header format", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	ipAddress := r.RemoteAddr
	issuer := "web"

	session, err := c.authService.Login(username, password, userAgent, ipAddress, issuer)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Write SessionID to a secure cookie (HttpOnly, Secure, SameSite) - In production, ensure Secure is true and SameSite is set appropriately
	http.SetCookie(w, &http.Cookie{
		Name:     "perlica_session",
		Domain:   c.config.GetHostDomain(),
		Path:     "/",
		Value:    session.SessionID,
		HttpOnly: true,
		Secure:   c.config.GetHostDomain() != "localhost:3000", // Set to true in production
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(session.SessionID))
}

// getSession retrieves the authenticated user's session info
func (c *AuthController) getSession(w http.ResponseWriter, r *http.Request) {
	session := guards.GetSessionFromContext(r)

	// If ?expanded=true is provided, include user info in the response
	expanded := r.URL.Query().Get("expanded") == "true"
	if expanded {
		user, err := c.userService.GetUserByID(session.UserID)
		if err != nil {
			http.Error(w, "Failed to retrieve user info", http.StatusInternalServerError)
			return
		}
		response := struct {
			Session any                      `json:"session"`
			User    *repositories.UserEntity `json:"user"`
		}{
			Session: map[string]any{
				"session_id":         session.SessionID,
				"session_expires_at": session.SessionExpiresAt,
			},
			User: user,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}
