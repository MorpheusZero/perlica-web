package repositories

import (
	"time"

	"github.com/morpheuszero/perlica-web/server/database"
)

type BotEntity struct {
	ID          int        `db:"id" json:"-"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	ModifiedAt  *time.Time `db:"modified_at" json:"modified_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
	ObjectID    string     `db:"object_id" json:"object_id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	DockerImage string     `db:"docker_image" json:"docker_image"`
}

type BotRepository struct {
	db *database.Database
}

func NewBotRepository(db *database.Database) *BotRepository {
	return &BotRepository{
		db: db,
	}
}

// CreateBot creates a new bot in the database.
func (r *BotRepository) CreateBot(name string, description *string, dockerImage string) (*BotEntity, error) {
	var bot BotEntity
	query := `
		INSERT INTO bots (name, description, docker_image)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.DB.QueryRow(query, name, description, dockerImage).Scan(&bot.ID)
	if err != nil {
		return nil, err
	}

	newBot, err := r.GetBotByID(bot.ID)
	if err != nil {
		return nil, err
	}
	return newBot, nil
}

func (r *BotRepository) GetBotByID(id int) (*BotEntity, error) {
	var bot BotEntity
	err := r.db.DB.Get(&bot, "SELECT * FROM bots WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

func (r *BotRepository) GetBotByObjectID(objectID string) (*BotEntity, error) {
	var bot BotEntity
	err := r.db.DB.Get(&bot, "SELECT * FROM bots WHERE object_id = $1 AND deleted_at IS NULL", objectID)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

func (r *BotRepository) DeleteBotByObjectID(objectID string) error {
	_, err := r.db.DB.Exec("UPDATE bots SET deleted_at = NOW() WHERE object_id = $1", objectID)
	return err
}
