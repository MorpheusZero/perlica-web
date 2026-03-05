package repositories

import (
	"time"

	"github.com/moprheuszero/perlica-web/server/database"
)

type SessionEntity struct {
	ID                 int        `db:"id" json:"-"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	ModifiedAt         *time.Time `db:"modified_at" json:"modified_at"`
	DeletedAt          *time.Time `db:"deleted_at" json:"deleted_at"`
	SessionID          string     `db:"session_id" json:"session_id"`
	SessionExpiresAt   time.Time  `db:"session_expires_at" json:"session_expires_at"`
	SessionMaxExpiryAt time.Time  `db:"session_max_expiry_at" json:"session_max_expiry_at"`
	UserID             int        `db:"user_id" json:"user_id"`
	UserAgent          string     `db:"user_agent" json:"user_agent"`
	IPAddress          string     `db:"ip_address" json:"ip_address"`
	Issuer             string     `db:"issuer" json:"issuer"`
}

type SessionRepository struct {
	db *database.Database
}

func NewSessionRepository(db *database.Database) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

// CreateSession creates a new session in the database.
func (r *SessionRepository) CreateSession(userID int, userAgent string, ipAddress string, issuer string) (*SessionEntity, error) {
	var session SessionEntity
	query := `
		INSERT INTO sessions (session_expires_at, session_max_expiry_at, user_id, user_agent, ip_address, issuer)
		VALUES (NOW() + INTERVAL '1 day', NOW() + INTERVAL '7 days', $1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.DB.QueryRow(query, userID, userAgent, ipAddress, issuer).Scan(&session.ID)
	if err != nil {
		return nil, err
	}

	newSession, err := r.GetSessionByID(session.ID)
	if err != nil {
		return nil, err
	}
	return newSession, nil
}

func (r *SessionRepository) GetSessionByID(id int) (*SessionEntity, error) {
	var session SessionEntity
	err := r.db.DB.Get(&session, "SELECT * FROM sessions WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSessionBySessionID retrieves a session by its session ID.
func (r *SessionRepository) GetSessionBySessionID(sessionID string) (*SessionEntity, error) {
	var session SessionEntity
	err := r.db.DB.Get(&session, "SELECT * FROM sessions WHERE session_id = $1 AND deleted_at IS NULL", sessionID)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSessionBySessionID marks a session as deleted in the database.
func (r *SessionRepository) DeleteSessionBySessionID(sessionID string) error {
	_, err := r.db.DB.Exec("UPDATE sessions SET deleted_at = NOW() WHERE session_id = $1", sessionID)
	return err
}

// ExtendSession extends the session expiration time by 1 day.
func (r *SessionRepository) ExtendSession(sessionID string) error {
	_, err := r.db.DB.Exec("UPDATE sessions SET session_expires_at = NOW() + INTERVAL '1 day' WHERE session_id = $1 AND deleted_at IS NULL", sessionID)
	return err
}
