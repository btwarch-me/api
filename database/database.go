package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	GitHubID    int64     `json:"github_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	AvatarURL   string    `json:"avatar_url"`
	AccessToken string    `json:"-"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

var DB *sql.DB

func Connect(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Connected to PostgreSQL database")
	return nil
}

func InitTables() error {
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		github_id BIGINT UNIQUE NOT NULL,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255),
		avatar_url TEXT,
		access_token TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_users_github_id ON users(github_id);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	`

	_, err := DB.Exec(createUsersTableSQL)
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	log.Println("Database tables initialized")
	return nil
}

func Close() error {
	return DB.Close()
}
