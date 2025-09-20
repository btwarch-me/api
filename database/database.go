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

type Record struct {
	ID                 uuid.UUID `json:"id"`
	UserId             uuid.UUID `json:"user_id"`
	RecordName         string    `json:"record_name"`
	RecordType         string    `json:"record_type"`
	RecordValue        string    `json:"record_value"`
	TTL                int       `json:"ttl"`
	IsActive           bool      `json:"is_active"`
	CloudflareRecordID *string   `json:"cloudflare_record_id"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
}

type SubdomainClaim struct {
	ID            uuid.UUID `json:"id"`
	UserId        uuid.UUID `json:"user_id"`
	SubdomainName string    `json:"subdomain_name"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
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
	migrationsDir := "./database/migrations"
	if err := RunMigrations(DB, migrationsDir); err != nil {
		return fmt.Errorf("error running migrations: %v", err)
	}

	log.Println("Database tables initialized via migrations")
	return nil
}

func Close() error {
	return DB.Close()
}
