package services

import (
	"github.com/moprheuszero/perlica-web/server/database/repositories"
	"github.com/moprheuszero/perlica-web/server/util"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

// GetUserByUsername retrieves a user by their username.
func (s *UserService) GetUserByUsername(username string) (*repositories.UserEntity, error) {
	return s.userRepository.GetUserByUsername(username)
}

// CreateUser creates a new user with the given username, user type, and password.
func (s *UserService) CreateUser(username string, userTypeKey string, password string) (*repositories.UserEntity, error) {
	passwordHash := util.HashPassword(password)
	return s.userRepository.CreateUser(username, userTypeKey, passwordHash)
}
