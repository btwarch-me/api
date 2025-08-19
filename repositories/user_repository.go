package repositories

import (
	"btwarch/database"
	"database/sql"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.DB}
}

func (r *UserRepository) CreateUser(githubID int64, username, email, avatarURL, accessToken string) (*database.User, error) {
	query := `
		INSERT INTO users (github_id, username, email, avatar_url, access_token)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, github_id, username, email, avatar_url, access_token, created_at, updated_at
	`

	user := &database.User{}
	err := r.db.QueryRow(
		query,
		githubID, username, email, avatarURL, accessToken,
	).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.Email,
		&user.AvatarURL, &user.AccessToken,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByGitHubID(githubID int64) (*database.User, error) {
	query := `
		SELECT id, github_id, username, email, avatar_url, access_token, created_at, updated_at
		FROM users WHERE github_id = $1
	`

	user := &database.User{}
	err := r.db.QueryRow(query, githubID).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.Email,
		&user.AvatarURL, &user.AccessToken,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	return user, nil
}

func (r *UserRepository) UpdateUserTokens(userID string, accessToken string) error {
	query := `
		UPDATE users 
		SET access_token = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, accessToken, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("error updating user tokens: %v", err)
	}

	return nil
}

func (r *UserRepository) InsertUser(user database.User) (int64, error) {
	res, err := r.db.Exec(
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

	return rowsAffected, nil
}
