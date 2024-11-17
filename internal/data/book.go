package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"github.com/lib/pq"
)

type Book struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Authors         []string  `json:"authors"`        // Multiple authors
	ISBN            string    `json:"isbn"`
	PublicationDate string    `json:"publication_date"`
	Genre           string    `json:"genre"`
	Description     string    `json:"description"`
	AverageRating   float64   `json:"average_rating"`
	CreatedAt       time.Time `json:"-"`
	Version         int32     `json:"version"`
}

type BookModel struct {
	DB *sql.DB
}

// Insert adds a new book to the database
func (m *BookModel) Insert(book *Book) error {
	query := `
		INSERT INTO books (title, authors, isbn, publication_date, genre, description, average_rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, version`

	// Convert authors to a string for the database
	authors := fmt.Sprintf("{%s}", joinStrings(book.Authors))

	args := []interface{}{
		book.Title,
		authors,
		book.ISBN,
		book.PublicationDate,
		book.Genre,
		book.Description,
		book.AverageRating,
	}

	return m.DB.QueryRowContext(context.Background(), query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

// Get retrieves a specific book by its ID
func (m *BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, title, authors, isbn, publication_date, genre, description, average_rating, created_at, version
		FROM books
		WHERE id = $1`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use pq.Array to scan the authors as a slice of strings
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		pq.Array(&book.Authors),  // Scan authors as a slice of strings
		&book.ISBN,
		&book.PublicationDate,
		&book.Genre,
		&book.Description,
		&book.AverageRating,
		&book.CreatedAt,
		&book.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil
}

func (m *BookModel) Update(book *Book) error {
	query := `
		UPDATE books
		SET title = $1, authors = $2, isbn = $3, publication_date = $4, genre = $5, description = $6, average_rating = $7, version = version + 1
		WHERE id = $8
		RETURNING version`

	// Use pq.Array to handle authors as a slice of strings
	args := []interface{}{
		book.Title,
		pq.Array(book.Authors),  // Use pq.Array for authors
		book.ISBN,
		book.PublicationDate,
		book.Genre,
		book.Description,
		book.AverageRating,
		book.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

// Delete deletes a book from the database by ID
func (m *BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM books
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// GetAll retrieves all books with optional filters and pagination
func (m *BookModel) GetAll(title, author, genre string, filters Filters) ([]*Book, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, title, authors, isbn, publication_date, genre, description, average_rating, created_at, version
		FROM books
		WHERE (title ILIKE $1 OR $1 = '')
		AND (authors ILIKE $2 OR $2 = '')  -- Added condition for authors
		AND (genre ILIKE $3 OR $3 = '')
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	// Joining authors if needed
	authors := "%" + author + "%" // For searching by author (ILIKE)

	args := []interface{}{
		"%" + title + "%",
		authors,
		"%" + genre + "%",
		filters.limit(),
		filters.offset(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	books := []*Book{}

	for rows.Next() {
		var book Book
		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.Title,
			&book.Authors,
			&book.ISBN,
			&book.PublicationDate,
			&book.Genre,
			&book.Description,
			&book.AverageRating,
			&book.CreatedAt,
			&book.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return books, metadata, nil
}

// Helper function to join a slice of strings into a single comma-separated string
func joinStrings(authors []string) string {
	return strings.Join(authors, ", ")
}
