-- Dropping the `created_by` column from the reading_lists table
ALTER TABLE reading_lists
DROP COLUMN created_by;

-- Dropping the `author_id` column from the reviews table
ALTER TABLE reviews
DROP COLUMN author_id;

-- Dropping the indexes created for the foreign key columns
DROP INDEX IF EXISTS idx_reading_lists_created_by;
DROP INDEX IF EXISTS idx_reviews_author_id;
