package repositories

import (
	"btwarch/config"
	"btwarch/database"
	"btwarch/services"
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

func (r *RecordRepository) getCloudflareService() (*services.CloudflareService, error) {
	cfg := config.LoadConfig()
	return services.NewCloudflareService(cfg.CloudFlareApiToken)
}

func (r *RecordRepository) CreateOnCloudflare(record database.Record) (string, error) {
	cf, err := r.getCloudflareService()
	if err != nil {
		return "", err
	}
	switch record.RecordType {
	case "A":
		resp, err := cf.AddARecord(record.RecordName, record.RecordValue)
		if err != nil {
			return "", err
		}
		return resp.ID, nil
	case "AAAA":
		resp, err := cf.AddAAAARecord(record.RecordName, record.RecordValue)
		if err != nil {
			return "", err
		}
		return resp.ID, nil
	case "TXT":
		resp, err := cf.AddTXTRecord(record.RecordName, record.RecordValue)
		if err != nil {
			return "", err
		}
		return resp.ID, nil
	case "CNAME":
		resp, err := cf.AddCNAMERecord(record.RecordName, record.RecordValue)
		if err != nil {
			return "", err
		}
		return resp.ID, nil
	default:
		return "", fmt.Errorf("invalid record type: %s", record.RecordType)
	}
}

func (r *RecordRepository) AddCloudflareRecord(record database.Record) (string, error) {
	return r.CreateOnCloudflare(record)
}

func (r *RecordRepository) UpdateCloudflareIDByNameAndType(recordName string, recordType string, cfID string) error {
	query := `
		UPDATE records
		SET cloudflare_record_id = $1, updated_at = $2
		WHERE record_name = $3 AND record_type = $4
	`
	_, err := r.db.Exec(query, cfID, time.Now(), recordName, recordType)
	return err
}

func (r *RecordRepository) CreateRecord(userID uuid.UUID, domainName, recordType, recordValue string, ttl int, isActive bool) (*database.Record, error) {
	var cloudflareID *string
	if isActive {
		newRecord := database.Record{
			UserId:      userID,
			RecordName:  domainName,
			RecordType:  recordType,
			RecordValue: recordValue,
			TTL:         ttl,
			IsActive:    true,
		}
		id, err := r.CreateOnCloudflare(newRecord)
		if err != nil {
			return nil, err
		}
		cloudflareID = &id
	}

	query := `
		INSERT INTO records (user_id, record_name, record_type, record_value, ttl, is_active, cloudflare_record_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, record_name, record_type, record_value, ttl, is_active, cloudflare_record_id, created_at, updated_at
	`

	record := &database.Record{}
	err := r.db.QueryRow(
		query,
		userID, domainName, recordType, recordValue, ttl, isActive, cloudflareID,
	).Scan(
		&record.ID, &record.UserId, &record.RecordName,
		&record.RecordType, &record.RecordValue, &record.TTL,
		&record.IsActive, &record.CloudflareRecordID, &record.CreatedAt, &record.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating record: %v", err)
	}

	return record, nil
}

func (r *RecordRepository) GetRecordsByUserID(userID uuid.UUID) ([]*database.Record, error) {
	query := `
		SELECT id, user_id, record_name, record_type, record_value, ttl, is_active, cloudflare_record_id, created_at, updated_at
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
			&record.ID, &record.UserId, &record.RecordName,
			&record.RecordType, &record.RecordValue, &record.TTL,
			&record.IsActive, &record.CloudflareRecordID, &record.CreatedAt, &record.UpdatedAt,
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
		SELECT id, user_id, record_name, record_type, record_value, ttl, is_active, cloudflare_record_id, created_at, updated_at
		FROM records WHERE id = $1
	`

	record := &database.Record{}
	err := r.db.QueryRow(query, recordID).Scan(
		&record.ID, &record.UserId, &record.RecordName,
		&record.RecordType, &record.RecordValue, &record.TTL,
		&record.IsActive, &record.CloudflareRecordID, &record.CreatedAt, &record.UpdatedAt,
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
		SELECT id, user_id, record_name, record_type, record_value, ttl, is_active, cloudflare_record_id, created_at, updated_at
		FROM records WHERE record_name = $1
	`

	record := &database.Record{}
	err := r.db.QueryRow(query, domainName).Scan(
		&record.ID, &record.UserId, &record.RecordName,
		&record.RecordType, &record.RecordValue, &record.TTL,
		&record.IsActive, &record.CloudflareRecordID, &record.CreatedAt, &record.UpdatedAt,
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
	// Load record to get Cloudflare ID and name/type
	rec, err := r.GetRecordByID(recordID)
	if err != nil {
		return err
	}
	if rec == nil {
		return nil
	}
	// Delete on Cloudflare if tracked
	if rec.CloudflareRecordID != nil && *rec.CloudflareRecordID != "" {
		cf, err := r.getCloudflareService()
		if err != nil {
			return err
		}
		if err := cf.DeleteRecordByID(*rec.CloudflareRecordID); err != nil {
			return fmt.Errorf("cloudflare delete failed: %w", err)
		}
	}

	query := `DELETE FROM records WHERE id = $1`
	_, err = r.db.Exec(query, recordID)
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
		`INSERT INTO records (user_id, record_name, record_type, record_value, ttl, is_active, cloudflare_record_id)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, record.RecordName, record.RecordType, record.RecordValue, record.TTL, record.IsActive, record.CloudflareRecordID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert record: %v", err)
	}

	return nil
}
