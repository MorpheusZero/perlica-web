package services

import (
	"errors"

	"github.com/moprheuszero/perlica-web/server/database/repositories"
	"github.com/moprheuszero/perlica-web/server/util"
)

type AuthService struct {
	sessionRepo *repositories.SessionRepository
	userService *UserService
}

func NewAuthService(sessionRepo *repositories.SessionRepository, userService *UserService) *AuthService {
	return &AuthService{
		sessionRepo: sessionRepo,
		userService: userService,
	}
}

func (s *AuthService) Login(username string, password string, userAgent string, ipAddress string, issuer string) (*repositories.SessionEntity, error) {
	user, err := s.userService.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if !util.VerifyPassword(password, *user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	session, err := s.sessionRepo.CreateSession(user.ID, userAgent, ipAddress, issuer)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves a session by session ID
func (s *AuthService) GetSession(sessionID string) (*repositories.SessionEntity, error) {
	return s.sessionRepo.GetSessionBySessionID(sessionID)
}
