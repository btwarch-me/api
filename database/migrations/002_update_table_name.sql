-- Migration: 002_update_table_name.sql
-- Description: Update domains table to records table with new fields

-- Rename domains table to records
ALTER TABLE domains RENAME TO records;

-- Drop existing indexes
DROP INDEX IF EXISTS idx_domains_user_id;
DROP INDEX IF EXISTS idx_domains_domain_name;

-- Drop cname_record column and add new columns
ALTER TABLE records 
    DROP COLUMN cname_record,
    ADD COLUMN record_type VARCHAR(10) NOT NULL,
    ADD COLUMN record_value TEXT NOT NULL,
    ADD COLUMN ttl INTEGER DEFAULT 1;

-- Create new indexes
CREATE INDEX IF NOT EXISTS idx_records_user_id ON records(user_id);
CREATE INDEX IF NOT EXISTS idx_records_domain_name ON records(domain_name);
CREATE INDEX IF NOT EXISTS idx_records_record_type ON records(record_type);
