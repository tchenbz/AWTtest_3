-- down_migration.sql

-- Drop indexes for reviews table
DROP INDEX IF EXISTS idx_reviews_book_id;
DROP INDEX IF EXISTS idx_reviews_rating;

-- Drop indexes for books table
DROP INDEX IF EXISTS idx_books_title;
DROP INDEX IF EXISTS idx_books_author;
DROP INDEX IF EXISTS idx_books_genre;

-- Drop indexes for reading lists table
DROP INDEX IF EXISTS idx_reading_lists_name;
DROP INDEX IF EXISTS idx_reading_lists_status;

-- Drop the 'reviews' table
DROP TABLE IF EXISTS reviews CASCADE;

-- Drop the 'reading_list_books' table (many-to-many relationship table)
DROP TABLE IF EXISTS reading_list_books CASCADE;

-- Drop the 'reading_lists' table
DROP TABLE IF EXISTS reading_lists CASCADE;

-- Drop the 'books' table
DROP TABLE IF EXISTS books CASCADE;

-- Drop the 'users' table
DROP TABLE IF EXISTS users CASCADE;
