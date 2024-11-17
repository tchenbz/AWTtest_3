package data

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"time"
// )

// // User represents a user in the system.
// type User struct {
//     ID            int64     `json:"id"`
//     Username      string    `json:"username"`
//     Email         string    `json:"email"`
//     Password      string    `json:"password"`
//     EmailVerified bool      `json:"email_verified"` // To store email verification status
//     VerificationToken string `json:"verification_token"` // Token for email verification
//     CreatedAt     time.Time `json:"created_at"`
//     Version       int32     `json:"version"`
// }


// // UserModel represents the model for interacting with users in the database.
// type UserModel struct {
// 	DB *sql.DB
// }

// // Insert inserts a new user into the database.
// func (m *UserModel) Insert(user *User) error {
// 	query := `
// 		INSERT INTO users (username, email, password)
// 		VALUES ($1, $2, $3)
// 		RETURNING id, created_at, version`

// 	args := []interface{}{user.Username, user.Email, user.Password}

// 	return m.DB.QueryRowContext(context.Background(), query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
// }

// // Get retrieves a specific user by ID.
// func (m *UserModel) Get(id int64) (*User, error) {
// 	if id < 1 {
// 		return nil, ErrRecordNotFound
// 	}

// 	query := `
// 		SELECT id, username, email, password, created_at, version
// 		FROM users
// 		WHERE id = $1`

// 	var user User

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	err := m.DB.QueryRowContext(ctx, query, id).Scan(
// 		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.Version,
// 	)

// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, ErrRecordNotFound
// 		default:
// 			return nil, err
// 		}
// 	}

// 	return &user, nil
// }

// // Update updates a user's details.
// func (m *UserModel) Update(user *User) error {
// 	query := `
// 		UPDATE users
// 		SET username = $1, email = $2, password = $3, version = version + 1
// 		WHERE id = $4
// 		RETURNING version`

// 	args := []interface{}{
// 		user.Username,
// 		user.Email,
// 		user.Password,
// 		user.ID,
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return ErrRecordNotFound
// 		default:
// 			return err
// 		}
// 	}

// 	return nil
// }

// // Delete removes a user from the database.
// func (m *UserModel) Delete(id int64) error {
// 	if id < 1 {
// 		return ErrRecordNotFound
// 	}

// 	query := `
// 		DELETE FROM users
// 		WHERE id = $1`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	result, err := m.DB.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return ErrRecordNotFound
// 	}

// 	return nil
// }

// // GetAll retrieves all users with optional filters and pagination.
// func (m *UserModel) GetAll(username, email string, filters Filters) ([]*User, Metadata, error) {
// 	query := fmt.Sprintf(`
// 		SELECT COUNT(*) OVER(), id, username, email, password, created_at, version
// 		FROM users
// 		WHERE (username ILIKE $1 OR $1 = '')
// 		AND (email ILIKE $2 OR $2 = '')
// 		ORDER BY %s %s, id ASC
// 		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

// 	args := []interface{}{
// 		"%" + username + "%",
// 		"%" + email + "%",
// 		filters.limit(),
// 		filters.offset(),
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	rows, err := m.DB.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, Metadata{}, err
// 	}
// 	defer rows.Close()

// 	totalRecords := 0
// 	users := []*User{}

// 	for rows.Next() {
// 		var user User
// 		err := rows.Scan(
// 			&totalRecords,
// 			&user.ID,
// 			&user.Username,
// 			&user.Email,
// 			&user.Password,
// 			&user.CreatedAt,
// 			&user.Version,
// 		)
// 		if err != nil {
// 			return nil, Metadata{}, err
// 		}
// 		users = append(users, &user)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, Metadata{}, err
// 	}

// 	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
// 	return users, metadata, nil
// }

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	//"golang.org/x/crypto/bcrypt"
	//"github.com/dgrijalva/jwt-go"
)

// User represents a user in the system.
type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	Version       int32     `json:"version"`
}

// UserModel represents the model for interacting with users in the database.
type UserModel struct {
	DB *sql.DB
}

// Insert a new user into the database.
func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, password, email_verified)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{user.Username, user.Email, user.Password, user.EmailVerified}

	return m.DB.QueryRowContext(context.Background(), query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
}

// Get a user by ID.
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

// GetByEmail retrieves a user by email.
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

// Update a user's details.
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

// Delete a user by ID.
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

// GetAll retrieves all users with optional filters and pagination.
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

