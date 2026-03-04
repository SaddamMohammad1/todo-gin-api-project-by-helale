-- Migration: Add user_id to todos table
ALTER TABLE todos ADD COLUMN user_id UUID NOT NULL;

-- Add foreign key constraint to link todos to users
ALTER TABLE todos ADD CONSTRAINT fk_todos_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE