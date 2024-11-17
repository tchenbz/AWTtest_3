-- Adding the `created_by` field to the reading_lists table to associate reading lists with users
ALTER TABLE reading_lists
ADD COLUMN created_by INT REFERENCES users(id) ON DELETE CASCADE;

-- Adding the `author_id` field to the reviews table to associate reviews with users
ALTER TABLE reviews
ADD COLUMN author_id INT REFERENCES users(id) ON DELETE CASCADE;

-- You can also add foreign key constraints or indexes if necessary for optimization
CREATE INDEX IF NOT EXISTS idx_reading_lists_created_by ON reading_lists(created_by);
CREATE INDEX IF NOT EXISTS idx_reviews_author_id ON reviews(author_id);
