package repositories

import (
	"btwarch/database"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type SubdomainClaimRepository struct {
	db *sql.DB
}

func NewSubdomainClaimRepository() *SubdomainClaimRepository {
	return &SubdomainClaimRepository{db: database.DB}
}

func (r *SubdomainClaimRepository) CreateClaim(userID uuid.UUID, subdomainName string) (*database.SubdomainClaim, error) {
	// Check if user already has a subdomain claim
	existingClaim, err := r.GetClaimByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing claims: %v", err)
	}
	if existingClaim != nil {
		return nil, fmt.Errorf("user already has a subdomain claim. Only one subdomain per user is allowed")
	}

	query := `
		INSERT INTO subdomain_claims (user_id, subdomain_name)
		VALUES ($1, $2)
		RETURNING id, user_id, subdomain_name, created_at, updated_at
	`

	claim := &database.SubdomainClaim{}
	err = r.db.QueryRow(
		query,
		userID, subdomainName,
	).Scan(
		&claim.ID, &claim.UserId, &claim.SubdomainName, &claim.CreatedAt, &claim.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating subdomain claim: %v", err)
	}

	return claim, nil
}

func (r *SubdomainClaimRepository) GetClaimBySubdomain(subdomainName string) (*database.SubdomainClaim, error) {
	query := `
		SELECT id, user_id, subdomain_name, created_at, updated_at
		FROM subdomain_claims WHERE subdomain_name = $1
	`

	claim := &database.SubdomainClaim{}
	err := r.db.QueryRow(query, subdomainName).Scan(
		&claim.ID, &claim.UserId, &claim.SubdomainName, &claim.CreatedAt, &claim.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting subdomain claim: %v", err)
	}

	return claim, nil
}

func (r *SubdomainClaimRepository) GetClaimByUserID(userID uuid.UUID) (*database.SubdomainClaim, error) {
	query := `
		SELECT id, user_id, subdomain_name, created_at, updated_at
		FROM subdomain_claims WHERE user_id = $1
	`

	claim := &database.SubdomainClaim{}
	err := r.db.QueryRow(query, userID).Scan(
		&claim.ID, &claim.UserId, &claim.SubdomainName, &claim.CreatedAt, &claim.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting subdomain claim: %v", err)
	}

	return claim, nil
}

func (r *SubdomainClaimRepository) GetClaimsByUserID(userID uuid.UUID) ([]*database.SubdomainClaim, error) {
	query := `
		SELECT id, user_id, subdomain_name, created_at, updated_at
		FROM subdomain_claims WHERE user_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting subdomain claims: %v", err)
	}
	defer rows.Close()

	var claims []*database.SubdomainClaim
	for rows.Next() {
		claim := &database.SubdomainClaim{}
		err := rows.Scan(
			&claim.ID, &claim.UserId, &claim.SubdomainName, &claim.CreatedAt, &claim.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning subdomain claim: %v", err)
		}
		claims = append(claims, claim)
	}

	return claims, nil
}

func (r *SubdomainClaimRepository) DeleteClaim(claimID uuid.UUID) error {
	query := `DELETE FROM subdomain_claims WHERE id = $1`
	_, err := r.db.Exec(query, claimID)
	if err != nil {
		return fmt.Errorf("error deleting subdomain claim: %v", err)
	}
	return nil
}
