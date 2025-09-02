package repositories

import (
	"btwarch/database"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MemeRepository struct {
	db *sql.DB
}

func NewMemeRepository() *MemeRepository {
	return &MemeRepository{db: database.DB}
}

func (r *MemeRepository) CreateMeme(userID uuid.UUID, title, description string, images []string) (*database.Meme, error) {
	query := `
		INSERT INTO memes (user_id, title, description, images)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, description, images, created_at, updated_at
	`

	meme := &database.Meme{}
	err := r.db.QueryRow(
		query,
		userID, title, description, images,
	).Scan(
		&meme.ID, &meme.UserId, &meme.Title, &meme.Description, &meme.Images,
		&meme.CreatedAt, &meme.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating meme: %v", err)
	}

	return meme, nil
}

func (r *MemeRepository) GetMemesByUserID(userID uuid.UUID) ([]*database.Meme, error) {
	query := `
		SELECT id, user_id, title, description, images, created_at, updated_at
		FROM memes WHERE user_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting memes: %v", err)
	}
	defer rows.Close()

	var memes []*database.Meme
	for rows.Next() {
		meme := &database.Meme{}
		err := rows.Scan(
			&meme.ID, &meme.UserId, &meme.Title, &meme.Description, &meme.Images,
			&meme.CreatedAt, &meme.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning meme: %v", err)
		}
		memes = append(memes, meme)
	}

	return memes, nil
}

func (r *MemeRepository) GetMemeByID(memeID uuid.UUID) (*database.Meme, error) {
	query := `
		SELECT id, user_id, title, description, images, created_at, updated_at
		FROM memes WHERE id = $1
	`

	meme := &database.Meme{}
	err := r.db.QueryRow(query, memeID).Scan(
		&meme.ID, &meme.UserId, &meme.Title, &meme.Description, &meme.Images,
		&meme.CreatedAt, &meme.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting meme: %v", err)
	}

	return meme, nil
}

func (r *MemeRepository) GetMemeByTitle(title string) (*database.Meme, error) {
	query := `
		SELECT id, user_id, title, description, images, created_at, updated_at
		FROM memes WHERE title = $1
	`

	meme := &database.Meme{}
	err := r.db.QueryRow(query, title).Scan(
		&meme.ID, &meme.UserId, &meme.Title, &meme.Description, &meme.Images,
		&meme.CreatedAt, &meme.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting meme: %v", err)
	}

	return meme, nil
}

func (r *MemeRepository) UpdateMeme(memeID uuid.UUID, title, description string, images []string) error {
	query := `
		UPDATE memes 
		SET updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, title, description, images, time.Now(), memeID)
	if err != nil {
		return fmt.Errorf("error updating meme: %v", err)
	}

	return nil
}

func (r *MemeRepository) DeleteMeme(memeID uuid.UUID) error {
	query := `DELETE FROM memes WHERE id = $1`

	_, err := r.db.Exec(query, memeID)
	if err != nil {
		return fmt.Errorf("error deleting meme: %v", err)
	}

	return nil
}
