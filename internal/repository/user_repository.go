package repository

import (
	"database/sql"

	"ticket-system/internal/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {

	var user model.User

	err := r.DB.QueryRow(`
		SELECT id, email, password_hash
		FROM users
		WHERE email = ?
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(
	email string,
	passwordHash string,
) (*model.User, error) {

	result, err := r.DB.Exec(`
		INSERT INTO users (email, password_hash)
		VALUES (?, ?)
	`, email, passwordHash)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}
