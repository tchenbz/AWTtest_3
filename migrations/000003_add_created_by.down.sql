-- Remove the 'created_by' field from the 'reading_lists' table
ALTER TABLE reading_lists
DROP COLUMN created_by;
