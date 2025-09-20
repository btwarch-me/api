-- Migration: 008_remove_things_related_to_memes.sql
-- Description: Remove memes table and related indexes

-- Drop indexes first
DROP INDEX IF EXISTS idx_memes_user_id;
DROP INDEX IF EXISTS idx_memes_title;

-- Drop memes table
DROP TABLE IF EXISTS memes CASCADE;
