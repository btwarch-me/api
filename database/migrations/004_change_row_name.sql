-- Migration: 004_change_row_name.sql
-- Description: Change row name from domain_name to record_name

ALTER TABLE records RENAME COLUMN domain_name TO record_name;   
