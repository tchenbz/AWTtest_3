package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ReadingList struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	Books       []int64   `json:"books"`      
	Status      string    `json:"status"`     
	CreatedAt   time.Time `json:"created_at"`
	Version     int32     `json:"version"`
}

type ReadingListModel struct {
	DB *sql.DB
}

func (m *ReadingListModel) Insert(readingList *ReadingList) error {
    query := `
        INSERT INTO reading_lists (name, description, created_by, status)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

    // Ensure the arguments match the columns in your table
    args := []interface{}{
        readingList.Name,
        readingList.Description,
        readingList.CreatedBy,
        readingList.Status,
    }

    // Insert into the database and return the inserted values
    err := m.DB.QueryRowContext(context.Background(), query, args...).Scan(
        &readingList.ID,        
        &readingList.CreatedAt, 
        &readingList.Version,  
    )

    return err
}


func (m *ReadingListModel) Get(id int64) (*ReadingList, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name, description, created_by, status, created_at, version
		FROM reading_lists
		WHERE id = $1`

	var readingList ReadingList

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&readingList.ID,
		&readingList.Name,
		&readingList.Description,
		&readingList.CreatedBy,
		&readingList.Status,
		&readingList.CreatedAt,
		&readingList.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &readingList, nil
}

func (m *ReadingListModel) Update(readingList *ReadingList) error {
	query := `
		UPDATE reading_lists
		SET name = $1, description = $2, created_by = $3, status = $4, version = version + 1
		WHERE id = $5
		RETURNING version`

	args := []interface{}{
		readingList.Name,
		readingList.Description,
		readingList.CreatedBy,
		readingList.Status,
		readingList.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&readingList.Version)
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

func (m *ReadingListModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM reading_lists
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

func (m *ReadingListModel) GetAll(name, status string, filters Filters) ([]*ReadingList, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, description, created_by, status, created_at, version
		FROM reading_lists
		WHERE (name ILIKE $1 OR $1 = '')
		AND (status ILIKE $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	args := []interface{}{
		"%" + name + "%",
		"%" + status + "%",
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
	readingLists := []*ReadingList{}

	for rows.Next() {
		var readingList ReadingList
		err := rows.Scan(
			&totalRecords,
			&readingList.ID,
			&readingList.Name,
			&readingList.Description,
			&readingList.CreatedBy,
			&readingList.Status,
			&readingList.CreatedAt,
			&readingList.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		readingLists = append(readingLists, &readingList)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return readingLists, metadata, nil
}

func (m *ReadingListModel) GetAllByUser(userID int64, filters Filters) ([]*ReadingList, Metadata, error) {
    // Construct the SQL query
    query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), id, name, description, created_by, status, created_at, version
        FROM reading_lists
        WHERE created_by = $1
        ORDER BY %s %s, id ASC
        LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

    // Arguments for the query
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
    readingLists := []*ReadingList{}

    for rows.Next() {
        var readingList ReadingList
        err := rows.Scan(
            &totalRecords,
            &readingList.ID,
            &readingList.Name,
            &readingList.Description,
            &readingList.CreatedBy,
            &readingList.Status,
            &readingList.CreatedAt,
            &readingList.Version,
        )
        if err != nil {
            return nil, Metadata{}, err
        }
        readingLists = append(readingLists, &readingList)
    }

    if err = rows.Err(); err != nil {
        return nil, Metadata{}, err
    }

    metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
    return readingLists, metadata, nil
}

