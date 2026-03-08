package guards

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/moprheuszero/perlica-web/server/database/repositories"
)

type contextKey string

const SessionContextKey contextKey = "perlica_session"

type AuthGuard struct {
	sessionRepo *repositories.SessionRepository
	userRepo    *repositories.UserRepository
}

func NewAuthGuard(sessionRepo *repositories.SessionRepository, userRepo *repositories.UserRepository) *AuthGuard {
	return &AuthGuard{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

// ValidateSession validates the session from perlica_session cookie and stores it in request context
func (g *AuthGuard) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract session ID from perlica_session cookie
		cookie, err := r.Cookie("perlica_session")
		if err != nil {
			fmt.Println("No perlica_session cookie found")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionID := cookie.Value

		// Retrieve and validate session
		session, err := g.sessionRepo.GetSessionBySessionID(sessionID)
		if err != nil || session == nil {
			fmt.Printf("Invalid session: %s\n", err.Error())
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Check if session has expired
		if time.Now().After(session.SessionExpiresAt) {
			fmt.Println("Session has expired")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Store session in context and proceed
		ctx := context.WithValue(r.Context(), SessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (g *AuthGuard) GetUserFromSessionContext(r *http.Request) *repositories.UserEntity {
	session := GetSessionFromContext(r)
	if session == nil {
		return nil
	}
	user, err := g.userRepo.GetUserByID(session.UserID)
	if err != nil {
		return nil
	}
	return user
}

// GetSessionFromContext retrieves the session from request context
func GetSessionFromContext(r *http.Request) *repositories.SessionEntity {
	session, ok := r.Context().Value(SessionContextKey).(*repositories.SessionEntity)
	if !ok {
		return nil
	}
	return session
}
