package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mini_project/internal/model"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, review *model.Review) error {
	const query = `
		INSERT INTO reviews (user_id, text, star_count)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, review.User.Id, review.Text, review.StarCount).Scan(&review.Id)
	if err != nil {
		return fmt.Errorf("create review: %w", err)
	}

	return nil
}

func (r *ReviewRepository) GetByID(ctx context.Context, id int) (*model.Review, error) {
	const query = `
		SELECT
			r.id,
			r.text,
			r.star_count,
			u.id,
			u.nickname
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		WHERE r.id = $1
	`

	var review model.Review
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.Id,
		&review.Text,
		&review.StarCount,
		&review.User.Id,
		&review.User.Nickname,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("review %d not found: %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("get review by id: %w", err)
	}

	return &review, nil
}

func (r *ReviewRepository) GetByUserID(ctx context.Context, userID int) ([]model.Review, error) {
	const query = `
		SELECT
			r.id,
			r.text,
			r.star_count,
			u.id,
			u.nickname
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		WHERE r.user_id = $1
		ORDER BY r.id
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get reviews by user id: %w", err)
	}
	defer rows.Close()

	return scanReviews(rows)
}

func (r *ReviewRepository) GetAll(ctx context.Context) ([]model.Review, error) {
	const query = `
		SELECT
			r.id,
			r.text,
			r.star_count,
			u.id,
			u.nickname
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		ORDER BY r.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all reviews: %w", err)
	}
	defer rows.Close()

	return scanReviews(rows)
}

func (r *ReviewRepository) Update(ctx context.Context, review *model.Review) error {
	const query = `
		UPDATE reviews
		SET user_id = $1, text = $2, star_count = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, review.User.Id, review.Text, review.StarCount, review.Id)
	if err != nil {
		return fmt.Errorf("update review: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update review rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("review %d not found", review.Id)
	}

	return nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id int) error {
	const query = `
		DELETE FROM reviews
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete review: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete review rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("review %d not found", id)
	}

	return nil
}

func scanReviews(rows *sql.Rows) ([]model.Review, error) {
	reviews := make([]model.Review, 0)
	for rows.Next() {
		var review model.Review
		if err := rows.Scan(
			&review.Id,
			&review.Text,
			&review.StarCount,
			&review.User.Id,
			&review.User.Nickname,
		); err != nil {
			return nil, fmt.Errorf("scan review: %w", err)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reviews: %w", err)
	}

	return reviews, nil
}
