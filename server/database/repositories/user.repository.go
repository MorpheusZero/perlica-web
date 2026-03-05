package repositories

import (
	"time"

	"github.com/moprheuszero/perlica-web/server/database"
)

type UserEntity struct {
	ID           int        `db:"id" json:"-"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	ModifiedAt   *time.Time `db:"modified_at" json:"modified_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at"`
	ObjectID     string     `db:"object_id" json:"object_id"`
	Username     string     `db:"username" json:"username"`
	UserTypeKey  string     `db:"user_type_key" json:"user_type_key"`
	PasswordHash *string    `db:"password_hash" json:"-"`
	APIKey       *string    `db:"api_key" json:"-"`
	LastLogin    *time.Time `db:"last_login" json:"last_login"`
}

type UserRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetUserByID retrieves a user by their ID.
func (r *UserRepository) GetUserByID(id int) (*UserEntity, error) {
	var user UserEntity
	err := r.db.DB.Get(&user, "SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func (r *UserRepository) GetUserByUsername(username string) (*UserEntity, error) {
	var user UserEntity
	err := r.db.DB.Get(&user, "SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(username string, userTypeKey string, passwordHash string) (*UserEntity, error) {
	var user UserEntity
	query := `
		INSERT INTO users (username, user_type_key, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.DB.QueryRow(query, username, userTypeKey, passwordHash).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
