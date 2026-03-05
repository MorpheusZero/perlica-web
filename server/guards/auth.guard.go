package guards

import (
	"context"
	"net/http"
	"time"

	"github.com/moprheuszero/perlica-web/server/database/repositories"
)

type contextKey string

const SessionContextKey contextKey = "perlica_session"

type AuthGuard struct {
	sessionRepo *repositories.SessionRepository
}

func NewAuthGuard(sessionRepo *repositories.SessionRepository) *AuthGuard {
	return &AuthGuard{
		sessionRepo: sessionRepo,
	}
}

// ValidateSession validates the session from perlica_session cookie and stores it in request context
func (g *AuthGuard) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract session ID from perlica_session cookie
		cookie, err := r.Cookie("perlica_session")
		if err != nil {
			http.Error(w, "Session cookie is required", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value

		// Retrieve and validate session
		session, err := g.sessionRepo.GetSessionBySessionID(sessionID)
		if err != nil || session == nil {
			http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
			return
		}

		// Check if session has expired
		if time.Now().After(session.SessionExpiresAt) {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		// Store session in context and proceed
		ctx := context.WithValue(r.Context(), SessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetSessionFromContext retrieves the session from request context
func GetSessionFromContext(r *http.Request) *repositories.SessionEntity {
	session, ok := r.Context().Value(SessionContextKey).(*repositories.SessionEntity)
	if !ok {
		return nil
	}
	return session
}
