package repositories

import (
	"btwarch/config"
	"btwarch/database"
	"btwarch/services"
	"database/sql"
	"fmt"
	"time"

	"github.com/cloudflare/cloudflare-go/v4/dns"
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

func (r *RecordRepository) CreateOnCloudflare(record database.Record) (*dns.RecordResponse, error) {
	cf, err := r.getCloudflareService()
	if err != nil {
		return nil, err
	}
	switch record.RecordType {
	case "A":
		resp, err := cf.AddARecord(record.RecordName, record.RecordValue)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case "AAAA":
		resp, err := cf.AddAAAARecord(record.RecordName, record.RecordValue)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case "TXT":
		resp, err := cf.AddTXTRecord(record.RecordName, record.RecordValue)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case "CNAME":
		resp, err := cf.AddCNAMERecord(record.RecordName, record.RecordValue)
		if err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, fmt.Errorf("invalid record type: %s", record.RecordType)
	}
}

func (r *RecordRepository) UpdateOnCloudflare(recordID string, record database.Record) (*dns.RecordResponse, error) {
	cf, err := r.getCloudflareService()
	if err != nil {
		return nil, err
	}
	switch record.RecordType {
	case "A":
		resp, err := cf.UpdateARecord(recordID, record.RecordName, record.RecordValue)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case "AAAA":
		resp, err := cf.UpdateAAAARecord(recordID, record.RecordName, record.RecordValue)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case "CNAME":
		resp, err := cf.UpdateCNAMERecord(recordID, record.RecordName, record.RecordValue)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case "TXT":
		resp, err := cf.UpdateTXTRecord(recordID, record.RecordName, record.RecordValue)
		if err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, fmt.Errorf("invalid record type: %s", record.RecordType)
	}
}

func (r *RecordRepository) CreateCloudflareRecord(record database.Record) (*dns.RecordResponse, error) {
	return r.CreateOnCloudflare(record)
}

func (r *RecordRepository) DeleteCloudflareRecord(recordID string) error {
	cf, err := r.getCloudflareService()
	if err != nil {
		return err
	}
	_, err = cf.DeleteRecordByID(recordID)
	return err
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
		id, err := r.CreateCloudflareRecord(newRecord)
		if err != nil {
			return nil, err
		}
		cloudflareID = &id.ID
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

func (r *RecordRepository) GetRecordByNameAndType(recordName string, recordType string) (*database.Record, error) {
	query := `
		SELECT id, user_id, record_name, record_type, record_value, ttl, is_active, cloudflare_record_id, created_at, updated_at
		FROM records WHERE record_name = $1 AND record_type = $2
	`

	record := &database.Record{}
	err := r.db.QueryRow(query, recordName, recordType).Scan(
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

func (r *RecordRepository) RecordExists(domainName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM records WHERE record_name = $1)`
	var exists bool
	err := r.db.QueryRow(query, domainName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking record: %v", err)
	}
	return exists, nil
}

func (r *RecordRepository) UpdateRecord(recordID uuid.UUID, recordName string, recordType string, recordValue string, ttl int) error {
	if recordType != "CNAME" && recordType != "A" && recordType != "AAAA" && recordType != "TXT" {
		return fmt.Errorf("invalid record type: %s", recordType)
	}

	existingRecord, err := r.GetRecordByID(recordID)
	if err != nil {
		return fmt.Errorf("error getting record: %v", err)
	}
	if existingRecord == nil {
		return fmt.Errorf("record not found")
	}

	if recordType != "TXT" && recordName != existingRecord.RecordName {
		return fmt.Errorf("cannot change record name for %s records", recordType)
	}

	query := `
		UPDATE records 
		SET record_name = $1, record_type = $2, record_value = $3, ttl = $4, updated_at = $5
		WHERE id = $6
	`

	_, err = r.db.Exec(query, recordName, recordType, recordValue, ttl, time.Now(), recordID)
	if err != nil {
		return fmt.Errorf("error updating record: %v", err)
	}
	return nil
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
	rec, err := r.GetRecordByID(recordID)
	if err != nil {
		return err
	}
	if rec == nil {
		return nil
	}
	if rec.CloudflareRecordID != nil && *rec.CloudflareRecordID != "" {
		cf, err := r.getCloudflareService()
		if err != nil {
			return err
		}
		if _, err := cf.DeleteRecordByID(*rec.CloudflareRecordID); err != nil {
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
