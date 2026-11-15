package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mini_project/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	const query = `
		INSERT INTO users (nickname)
		VALUES ($1)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, user.Nickname).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	const query = `
		SELECT id, nickname
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Nickname)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %d not found: %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	const query = `
		SELECT id, nickname
		FROM users
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Nickname); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	const query = `
		UPDATE users
		SET nickname = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, user.Nickname, user.Id)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update user rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user %d not found", user.Id)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	const query = `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user %d not found", id)
	}

	return nil
}
