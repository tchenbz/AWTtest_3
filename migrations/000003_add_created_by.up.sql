-- Add the 'created_by' field to the 'reading_lists' table to associate it with a user
ALTER TABLE reading_lists
ADD COLUMN created_by INT REFERENCES users(id) ON DELETE CASCADE;
