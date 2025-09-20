-- Migration: 007_create_subdomain_claims.sql
-- Description: Create subdomain claims table for managing user subdomain ownership

-- Create subdomain_claims table
CREATE TABLE IF NOT EXISTS subdomain_claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subdomain_name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for subdomain_claims table
CREATE INDEX IF NOT EXISTS idx_subdomain_claims_user_id ON subdomain_claims(user_id);
CREATE INDEX IF NOT EXISTS idx_subdomain_claims_subdomain_name ON subdomain_claims(subdomain_name);
