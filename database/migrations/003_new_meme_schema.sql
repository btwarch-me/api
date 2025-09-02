-- Migration: 003_new_meme_schema.sql
-- Description: Database schema for memes table

-- Create memes table
CREATE TABLE IF NOT EXISTS memes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    images TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for memes table
CREATE INDEX IF NOT EXISTS idx_memes_user_id ON memes(user_id);
CREATE INDEX IF NOT EXISTS idx_memes_title ON memes(title);
