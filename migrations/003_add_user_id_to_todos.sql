-- Add user_id column to todos table for ownership
ALTER TABLE todos ADD COLUMN IF NOT EXISTS user_id UUID;

-- Add foreign key constraint
ALTER TABLE todos 
    ADD CONSTRAINT fk_todos_user_id 
    FOREIGN KEY (user_id) 
    REFERENCES users(id) 
    ON DELETE CASCADE;

-- Create index on user_id for faster queries
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- Make user_id NOT NULL for new records (optional, based on requirements)
-- ALTER TABLE todos ALTER COLUMN user_id SET NOT NULL;
