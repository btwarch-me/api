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

type Domain struct {
	ID         uuid.UUID `json:"id"`
	UserId     uuid.UUID `json:"user_id"`
	DomainName string    `json:"domain_name"`
	Cname      string    `json:"cname_record"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
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

	CREATE TABLE IF NOT EXISTS domains (
    	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   	 	domain_name VARCHAR(255) UNIQUE NOT NULL,
    	cname_record VARCHAR(255) NOT NULL,
    	is_active BOOLEAN DEFAULT TRUE,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

	CREATE INDEX IF NOT EXISTS idx_domains_user_id ON domains(user_id);
	CREATE INDEX IF NOT EXISTS idx_domains_domain_name ON domains(domain_name);
	`

	_, err := DB.Exec(createUsersTableSQL)
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	log.Println("Database tables initialized")
	return nil
}

func InsertUser(user User) (int64, error) {
	res, err := DB.Exec(
		`INSERT INTO users (github_id, username, email, avatar_url, access_token)
         VALUES ($1, $2, $3, $4, $5)`,
		user.GitHubID, user.Username, user.Email, user.AvatarURL, user.AccessToken,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("no user was inserted")
	}

	log.Printf("User inserted successfully: GitHubID=%d, RowsAffected=%d\n", user.GitHubID, rowsAffected)
	return rowsAffected, nil
}

func AddDomain(githubID int64, domain Domain) error {
	var userID uuid.UUID
	err := DB.QueryRow(`SELECT id FROM users WHERE github_id=$1`, githubID).Scan(&userID)
	if err != nil {
		return fmt.Errorf("user not found or query failed: %v", err)
	}

	_, err = DB.Exec(
		`INSERT INTO domains (user_id, domain_name, cname_record, is_active)
         VALUES ($1, $2, $3, $4)`,
		userID, domain.DomainName, domain.Cname, domain.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to insert domain: %v", err)
	}

	log.Printf("Domain inserted successfully: %s for user %s\n", domain.DomainName, userID.String())
	return nil
}

func Close() error {
	return DB.Close()
}