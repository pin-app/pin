package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/models"
)

type ratingRepository struct {
	db *database.DB
}

func NewRatingRepository(db *database.DB) RatingRepository {
	return &ratingRepository{db: db}
}

func (r *ratingRepository) CreateRating(ctx context.Context, rating *models.PlaceRating) error {
	query := `
		INSERT INTO place_ratings (id, user_id, place_id, rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, place_id) 
		DO UPDATE SET rating = EXCLUDED.rating, updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		rating.ID, rating.UserID, rating.PlaceID, rating.Rating, rating.CreatedAt, rating.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create place rating: %w", err)
	}

	return nil
}

func (r *ratingRepository) GetRating(ctx context.Context, userID, placeID uuid.UUID) (*models.PlaceRating, error) {
	query := `
		SELECT id, user_id, place_id, rating, created_at, updated_at, deleted_at
		FROM place_ratings
		WHERE user_id = $1 AND place_id = $2 AND deleted_at IS NULL
	`

	rating := &models.PlaceRating{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, userID, placeID).Scan(
		&rating.ID, &rating.UserID, &rating.PlaceID, &rating.Rating,
		&rating.CreatedAt, &rating.UpdatedAt, &rating.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("place rating not found")
		}
		return nil, fmt.Errorf("failed to get place rating: %w", err)
	}

	return rating, nil
}

func (r *ratingRepository) UpdateRating(ctx context.Context, rating *models.PlaceRating) error {
	query := `
		UPDATE place_ratings
		SET rating = $3, updated_at = $4
		WHERE user_id = $1 AND place_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query,
		rating.UserID, rating.PlaceID, rating.Rating, rating.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update place rating: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("place rating not found or already deleted")
	}

	return nil
}

func (r *ratingRepository) DeleteRating(ctx context.Context, userID, placeID uuid.UUID) error {
	query := `
		UPDATE place_ratings
		SET deleted_at = NOW()
		WHERE user_id = $1 AND place_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, userID, placeID)
	if err != nil {
		return fmt.Errorf("failed to delete place rating: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("place rating not found or already deleted")
	}

	return nil
}

func (r *ratingRepository) GetRatingsByPlaceID(ctx context.Context, placeID uuid.UUID, limit, offset int) ([]*models.PlaceRating, error) {
	query := `
		SELECT id, user_id, place_id, rating, created_at, updated_at, deleted_at
		FROM place_ratings
		WHERE place_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, placeID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get ratings by place ID: %w", err)
	}
	defer rows.Close()

	var ratings []*models.PlaceRating
	for rows.Next() {
		rating := &models.PlaceRating{}
		err := rows.Scan(
			&rating.ID, &rating.UserID, &rating.PlaceID, &rating.Rating,
			&rating.CreatedAt, &rating.UpdatedAt, &rating.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan place rating: %w", err)
		}
		ratings = append(ratings, rating)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate place ratings: %w", err)
	}

	return ratings, nil
}

func (r *ratingRepository) GetRatingsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.PlaceRating, error) {
	query := `
		SELECT id, user_id, place_id, rating, created_at, updated_at, deleted_at
		FROM place_ratings
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get ratings by user ID: %w", err)
	}
	defer rows.Close()

	var ratings []*models.PlaceRating
	for rows.Next() {
		rating := &models.PlaceRating{}
		err := rows.Scan(
			&rating.ID, &rating.UserID, &rating.PlaceID, &rating.Rating,
			&rating.CreatedAt, &rating.UpdatedAt, &rating.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan place rating: %w", err)
		}
		ratings = append(ratings, rating)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate place ratings: %w", err)
	}

	return ratings, nil
}

func (r *ratingRepository) GetAverageRating(ctx context.Context, placeID uuid.UUID) (float64, int, error) {
	query := `
		SELECT COALESCE(AVG(rating), 0), COUNT(*)
		FROM place_ratings
		WHERE place_id = $1 AND deleted_at IS NULL
	`

	var avgRating float64
	var count int
	err := r.db.GetConnection().QueryRowContext(ctx, query, placeID).Scan(&avgRating, &count)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return avgRating, count, nil
}

func (r *ratingRepository) CreateComparison(ctx context.Context, comparison *models.PlaceComparison) error {
	query := `
		INSERT INTO place_comparisons (id, user_id, better_place_id, worse_place_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, better_place_id, worse_place_id) 
		DO UPDATE SET updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		comparison.ID, comparison.UserID, comparison.BetterPlaceID, comparison.WorsePlaceID,
		comparison.CreatedAt, comparison.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create place comparison: %w", err)
	}

	return nil
}

func (r *ratingRepository) GetComparison(ctx context.Context, userID, betterPlaceID, worsePlaceID uuid.UUID) (*models.PlaceComparison, error) {
	query := `
		SELECT id, user_id, better_place_id, worse_place_id, created_at, updated_at, deleted_at
		FROM place_comparisons
		WHERE user_id = $1 AND better_place_id = $2 AND worse_place_id = $3 AND deleted_at IS NULL
	`

	comparison := &models.PlaceComparison{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, userID, betterPlaceID, worsePlaceID).Scan(
		&comparison.ID, &comparison.UserID, &comparison.BetterPlaceID, &comparison.WorsePlaceID,
		&comparison.CreatedAt, &comparison.UpdatedAt, &comparison.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("place comparison not found")
		}
		return nil, fmt.Errorf("failed to get place comparison: %w", err)
	}

	return comparison, nil
}

func (r *ratingRepository) DeleteComparison(ctx context.Context, userID, betterPlaceID, worsePlaceID uuid.UUID) error {
	query := `
		UPDATE place_comparisons
		SET deleted_at = NOW()
		WHERE user_id = $1 AND better_place_id = $2 AND worse_place_id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, userID, betterPlaceID, worsePlaceID)
	if err != nil {
		return fmt.Errorf("failed to delete place comparison: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("place comparison not found or already deleted")
	}

	return nil
}

func (r *ratingRepository) GetComparisonsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.PlaceComparison, error) {
	query := `
		SELECT id, user_id, better_place_id, worse_place_id, created_at, updated_at, deleted_at
		FROM place_comparisons
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comparisons by user ID: %w", err)
	}
	defer rows.Close()

	var comparisons []*models.PlaceComparison
	for rows.Next() {
		comparison := &models.PlaceComparison{}
		err := rows.Scan(
			&comparison.ID, &comparison.UserID, &comparison.BetterPlaceID, &comparison.WorsePlaceID,
			&comparison.CreatedAt, &comparison.UpdatedAt, &comparison.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan place comparison: %w", err)
		}
		comparisons = append(comparisons, comparison)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate place comparisons: %w", err)
	}

	return comparisons, nil
}
