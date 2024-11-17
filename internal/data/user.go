package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	Version       int32     `json:"version"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, password, email_verified)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{user.Username, user.Email, user.Password, user.EmailVerified}

	return m.DB.QueryRowContext(context.Background(), query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
}

func (m *UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, username, email, password, email_verified, created_at, version
		FROM users
		WHERE id = $1`

	var user User
	err := m.DB.QueryRowContext(context.Background(), query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.EmailVerified, &user.CreatedAt, &user.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, username, email, password, email_verified, created_at, version
		FROM users
		WHERE email = $1`

	var user User
	err := m.DB.QueryRowContext(context.Background(), query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.EmailVerified, &user.CreatedAt, &user.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, email_verified = $4, version = version + 1
		WHERE id = $5
		RETURNING version`

	args := []interface{}{user.Username, user.Email, user.Password, user.EmailVerified, user.ID}

	err := m.DB.QueryRowContext(context.Background(), query, args...).Scan(&user.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (m *UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM users
		WHERE id = $1`

	result, err := m.DB.ExecContext(context.Background(), query, id)
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

func (m *UserModel) GetAll(username, email string, filters Filters) ([]*User, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, username, email, password, email_verified, created_at, version
		FROM users
		WHERE (username ILIKE $1 OR $1 = '')
		AND (email ILIKE $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	args := []interface{}{
		"%" + username + "%",
		"%" + email + "%",
		filters.limit(),
		filters.offset(),
	}

	rows, err := m.DB.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	users := []*User{}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.EmailVerified,
			&user.CreatedAt,
			&user.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return users, metadata, nil
}

