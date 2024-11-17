package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type Review struct {
	ID           int64     `json:"id"`
	BookID       int64     `json:"book_id"`        // Changed from product_id to book_id
	Content      string    `json:"content"`
	Author       string    `json:"author"`
	Rating       int       `json:"rating"`         
	HelpfulCount int       `json:"helpful_count"`  
	CreatedAt    time.Time `json:"created_at"`
	Version      int32     `json:"version"`
}

type ReviewModel struct {
	DB *sql.DB
}

// Insert adds a new review to the database
func (m ReviewModel) Insert(review *Review) error {
	query := `
		INSERT INTO reviews (book_id, content, author, rating, helpful_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`

	args := []interface{}{review.BookID, review.Content, review.Author, review.Rating, review.HelpfulCount}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID, &review.CreatedAt, &review.Version)
}

// Get retrieves a specific review for a book by review ID
func (m ReviewModel) Get(bookID, reviewID int64) (*Review, error) {
	if bookID < 1 || reviewID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, book_id, content, author, rating, helpful_count, created_at, version
		FROM reviews
		WHERE book_id = $1 AND id = $2`

	var review Review

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, bookID, reviewID).Scan(
		&review.ID, &review.BookID, &review.Content, &review.Author, 
		&review.Rating, &review.HelpfulCount, &review.CreatedAt, &review.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &review, nil
}

// Update modifies an existing review
func (m ReviewModel) Update(review *Review) error {
	query := `
		UPDATE reviews
		SET content = $1, author = $2, rating = $3, helpful_count = $4, version = version + 1
		WHERE book_id = $5 AND id = $6
		RETURNING version`

	args := []interface{}{review.Content, review.Author, review.Rating, review.HelpfulCount, review.BookID, review.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&review.Version)
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

// Delete removes a review for a book from the database
func (m ReviewModel) Delete(bookID, reviewID int64) error {
	if bookID < 1 || reviewID < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM reviews
		WHERE book_id = $1 AND id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, bookID, reviewID)
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

// GetAll retrieves all reviews, with optional filters
func (m ReviewModel) GetAll(content, author string, rating int, filters Filters) ([]*Review, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, book_id, content, author, rating, helpful_count, created_at, version
		FROM reviews
		WHERE (content ILIKE $1 OR $1 = '')
		AND (author ILIKE $2 OR $2 = '')
		AND (rating = $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	args := []interface{}{
		"%" + content + "%",
		"%" + author + "%",
		rating,
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
	reviews := []*Review{}

	for rows.Next() {
		var review Review
		err := rows.Scan(
			&totalRecords,
			&review.ID,
			&review.BookID,
			&review.Content,
			&review.Author,
			&review.Rating,
			&review.HelpfulCount,
			&review.CreatedAt,
			&review.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return reviews, metadata, nil
}

// GetAllForBook retrieves reviews for a specific book
func (m ReviewModel) GetAllForBook(bookID int64, content, author string, rating int, filters Filters) ([]*Review, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, book_id, content, author, rating, helpful_count, created_at, version
		FROM reviews
		WHERE book_id = $1
		AND (content ILIKE $2 OR $2 = '')
		AND (author ILIKE $3 OR $3 = '')
		AND (rating = $4 OR $4 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	args := []interface{}{
		bookID,
		"%" + content + "%",
		"%" + author + "%",
		rating,
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
	reviews := []*Review{}

	for rows.Next() {
		var review Review
		err := rows.Scan(
			&totalRecords,
			&review.ID,
			&review.BookID,
			&review.Content,
			&review.Author,
			&review.Rating,
			&review.HelpfulCount,
			&review.CreatedAt,
			&review.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return reviews, metadata, nil
}

// GetAllByUser retrieves all reviews written by a user
func (m *ReviewModel) GetAllByUser(userID int64, filters Filters) ([]*Review, Metadata, error) {
	// // Validate filters
	// v := validator.New()
	// ValidateFilters(v, filters)
	// if !v.IsEmpty() {
	// 	return nil, Metadata{}, fmt.Errorf("invalid filters: %v", v.Errors)
	// }

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, book_id, content, author, rating, helpful_count, created_at, version
		FROM reviews
		WHERE author = $1
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	args := []interface{}{
		userID,
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
	reviews := []*Review{}

	for rows.Next() {
		var review Review
		err := rows.Scan(
			&totalRecords,
			&review.ID,
			&review.BookID,
			&review.Content,
			&review.Author,
			&review.Rating,
			&review.HelpfulCount,
			&review.CreatedAt,
			&review.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return reviews, metadata, nil
}
