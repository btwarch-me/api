-- Migration: 005_add_new_row.sql
-- Description: Add new row to records table

ALTER TABLE records ADD COLUMN cloudflare_record_id VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_records_cloudflare_record_id ON records(cloudflare_record_id);

