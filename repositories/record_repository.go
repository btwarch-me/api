package repositories

import (
	"btwarch/database"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RecordRepository struct {
	db *sql.DB
}

func NewRecordRepository() *RecordRepository {
	return &RecordRepository{db: database.DB}
}

func (r *RecordRepository) CreateRecord(userID uuid.UUID, domainName, recordType, recordValue string, ttl int, isActive bool) (*database.Record, error) {
	query := `
		INSERT INTO records (user_id, domain_name, record_type, record_value, ttl, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, domain_name, record_type, record_value, ttl, is_active, created_at, updated_at
	`

	record := &database.Record{}
	err := r.db.QueryRow(
		query,
		userID, domainName, recordType, recordValue, ttl, isActive,
	).Scan(
		&record.ID, &record.UserId, &record.DomainName,
		&record.RecordType, &record.RecordValue, &record.TTL,
		&record.IsActive, &record.CreatedAt, &record.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating record: %v", err)
	}

	return record, nil
}

func (r *RecordRepository) GetRecordsByUserID(userID uuid.UUID) ([]*database.Record, error) {
	query := `
		SELECT id, user_id, domain_name, record_type, record_value, ttl, is_active, created_at, updated_at
		FROM records WHERE user_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting records: %v", err)
	}
	defer rows.Close()

	var records []*database.Record
	for rows.Next() {
		record := &database.Record{}
		err := rows.Scan(
			&record.ID, &record.UserId, &record.DomainName,
			&record.RecordType, &record.RecordValue, &record.TTL,
			&record.IsActive, &record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning record: %v", err)
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *RecordRepository) GetRecordByID(recordID uuid.UUID) (*database.Record, error) {
	query := `
		SELECT id, user_id, domain_name, record_type, record_value, ttl, is_active, created_at, updated_at
		FROM records WHERE id = $1
	`

	record := &database.Record{}
	err := r.db.QueryRow(query, recordID).Scan(
		&record.ID, &record.UserId, &record.DomainName,
		&record.RecordType, &record.RecordValue, &record.TTL,
		&record.IsActive, &record.CreatedAt, &record.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting record: %v", err)
	}

	return record, nil
}

func (r *RecordRepository) GetRecordByName(domainName string) (*database.Record, error) {
	query := `
		SELECT id, user_id, domain_name, record_type, record_value, ttl, is_active, created_at, updated_at
		FROM records WHERE domain_name = $1
	`

	record := &database.Record{}
	err := r.db.QueryRow(query, domainName).Scan(
		&record.ID, &record.UserId, &record.DomainName,
		&record.RecordType, &record.RecordValue, &record.TTL,
		&record.IsActive, &record.CreatedAt, &record.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting record: %v", err)
	}

	return record, nil
}

func (r *RecordRepository) UpdateRecordStatus(recordID uuid.UUID, isActive bool) error {
	query := `
		UPDATE records 
		SET is_active = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, isActive, time.Now(), recordID)
	if err != nil {
		return fmt.Errorf("error updating record status: %v", err)
	}

	return nil
}

func (r *RecordRepository) DeleteRecord(recordID uuid.UUID) error {
	query := `DELETE FROM records WHERE id = $1`

	_, err := r.db.Exec(query, recordID)
	if err != nil {
		return fmt.Errorf("error deleting record: %v", err)
	}

	return nil
}

func (r *RecordRepository) AddRecordByGitHubID(githubID int64, record database.Record) error {

	var userID uuid.UUID
	err := r.db.QueryRow(`SELECT id FROM users WHERE github_id = $1`, githubID).Scan(&userID)
	if err != nil {
		return fmt.Errorf("user not found or query failed: %v", err)
	}

	_, err = r.db.Exec(
		`INSERT INTO records (user_id, domain_name, record_type, record_value, ttl, is_active)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, record.DomainName, record.RecordType, record.RecordValue, record.TTL, record.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to insert record: %v", err)
	}

	return nil
}
