package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sqlx.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Initialize(connectionString string) error {

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return err
	}

	d.DB = db
	fmt.Println("Database connection established successfully")
	return nil
}
