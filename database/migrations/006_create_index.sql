-- Migration: 006_change_index_name.sql
-- Description: Create index on record_name

DROP INDEX IF EXISTS idx_records_domain_name;

CREATE INDEX IF NOT EXISTS idx_records_record_name ON records(record_name);