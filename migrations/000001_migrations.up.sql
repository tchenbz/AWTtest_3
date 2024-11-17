-- migration.sql

-- Create the 'users' table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version INT DEFAULT 1
);

-- Create the 'books' table
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    authors TEXT[] NOT NULL,
    isbn VARCHAR(20),
    publication_date DATE,
    genre VARCHAR(100),
    description TEXT,
    average_rating FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version INT DEFAULT 1
);

-- Create the 'reading_lists' table
CREATE TABLE IF NOT EXISTS reading_lists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INT NOT NULL,  -- Foreign Key to users
    status VARCHAR(50) CHECK (status IN ('currently reading', 'completed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version INT DEFAULT 1,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Create the 'reading_list_books' table to associate books with reading lists (many-to-many relationship)
CREATE TABLE IF NOT EXISTS reading_list_books (
    reading_list_id INT NOT NULL,
    book_id INT NOT NULL,
    PRIMARY KEY (reading_list_id, book_id),
    FOREIGN KEY (reading_list_id) REFERENCES reading_lists(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- Create the 'reviews' table
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,  -- Foreign Key to books
    content TEXT,
    author VARCHAR(255),
    rating INT CHECK (rating >= 1 AND rating <= 5),
    helpful_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version INT DEFAULT 1,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- Indexing the 'reviews' table for better query performance on book_id and rating
CREATE INDEX IF NOT EXISTS idx_reviews_book_id ON reviews(book_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);

-- Create indexes on frequently queried columns in the books table
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
CREATE INDEX IF NOT EXISTS idx_books_author ON books(authors);
CREATE INDEX IF NOT EXISTS idx_books_genre ON books(genre);

-- Create indexes for efficient searching and filtering on reading lists
CREATE INDEX IF NOT EXISTS idx_reading_lists_name ON reading_lists(name);
CREATE INDEX IF NOT EXISTS idx_reading_lists_status ON reading_lists(status);
