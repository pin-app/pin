package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/models"
)

type placeRepository struct {
	db *database.DB
}

func NewPlaceRepository(db *database.DB) PlaceRepository {
	return &placeRepository{db: db}
}

func (r *placeRepository) Create(ctx context.Context, place *models.Place) error {
	query := `
		INSERT INTO places (id, name, geometry, properties, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	// Convert Properties to JSON
	var propertiesJSON []byte
	var err error
	if place.Properties != nil {
		propertiesJSON, err = json.Marshal(place.Properties)
		if err != nil {
			return fmt.Errorf("failed to marshal properties: %w", err)
		}
	}

	_, err = r.db.GetConnection().ExecContext(ctx, query,
		place.ID, place.Name, place.Geometry, propertiesJSON, place.CreatedAt, place.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create place: %w", err)
	}

	return nil
}

func (r *placeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Place, error) {
	query := `
		SELECT id, name, geometry, properties, created_at, updated_at, deleted_at
		FROM places
		WHERE id = $1 AND deleted_at IS NULL
	`

	place := &models.Place{}
	var propertiesJSON []byte
	err := r.db.GetConnection().QueryRowContext(ctx, query, id).Scan(
		&place.ID, &place.Name, &place.Geometry, &propertiesJSON,
		&place.CreatedAt, &place.UpdatedAt, &place.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("place not found")
		}
		return nil, fmt.Errorf("failed to get place by ID: %w", err)
	}

	// Unmarshal properties JSON
	if len(propertiesJSON) > 0 {
		err = json.Unmarshal(propertiesJSON, &place.Properties)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal properties: %w", err)
		}
	}

	return place, nil
}

func (r *placeRepository) Update(ctx context.Context, place *models.Place) error {
	query := `
		UPDATE places
		SET name = $2, geometry = $3, properties = $4, updated_at = $5
		WHERE id = $1 AND deleted_at IS NULL
	`

	// Convert Properties to JSON
	var propertiesJSON []byte
	var err error
	if place.Properties != nil {
		propertiesJSON, err = json.Marshal(place.Properties)
		if err != nil {
			return fmt.Errorf("failed to marshal properties: %w", err)
		}
	}

	result, err := r.db.GetConnection().ExecContext(ctx, query,
		place.ID, place.Name, place.Geometry, propertiesJSON, place.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update place: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("place not found or already deleted")
	}

	return nil
}

func (r *placeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE places
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete place: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("place not found or already deleted")
	}

	return nil
}

func (r *placeRepository) List(ctx context.Context, limit, offset int) ([]*models.Place, error) {
	query := `
		SELECT id, name, geometry, properties, created_at, updated_at, deleted_at
		FROM places
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}
	defer rows.Close()

	var places []*models.Place
	for rows.Next() {
		place := &models.Place{}
		var propertiesJSON []byte
		err := rows.Scan(
			&place.ID, &place.Name, &place.Geometry, &propertiesJSON,
			&place.CreatedAt, &place.UpdatedAt, &place.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan place: %w", err)
		}

		// Unmarshal properties JSON
		if len(propertiesJSON) > 0 {
			err = json.Unmarshal(propertiesJSON, &place.Properties)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal properties: %w", err)
			}
		}

		places = append(places, place)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate places: %w", err)
	}

	return places, nil
}

func (r *placeRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.Place, error) {
	searchQuery := `
		SELECT id, name, geometry, properties, created_at, updated_at, deleted_at
		FROM places
		WHERE deleted_at IS NULL
		AND LOWER(name) LIKE LOWER($1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	searchTerm := "%" + strings.ToLower(query) + "%"
	rows, err := r.db.GetConnection().QueryContext(ctx, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search places: %w", err)
	}
	defer rows.Close()

	var places []*models.Place
	for rows.Next() {
		place := &models.Place{}
		var propertiesJSON []byte
		err := rows.Scan(
			&place.ID, &place.Name, &place.Geometry, &propertiesJSON,
			&place.CreatedAt, &place.UpdatedAt, &place.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan place: %w", err)
		}

		// Unmarshal properties JSON
		if len(propertiesJSON) > 0 {
			err = json.Unmarshal(propertiesJSON, &place.Properties)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal properties: %w", err)
			}
		}

		places = append(places, place)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate places: %w", err)
	}

	return places, nil
}

func (r *placeRepository) SearchNearby(ctx context.Context, lat, lng float64, radiusKm float64, limit int) ([]*models.Place, error) {
	query := `
		SELECT id, name, geometry, properties, created_at, updated_at, deleted_at,
			ST_Distance(geometry, ST_SetSRID(ST_MakePoint($2, $1), 4326)) as distance
		FROM places
		WHERE deleted_at IS NULL
		AND ST_DWithin(geometry, ST_SetSRID(ST_MakePoint($2, $1), 4326), $3)
		ORDER BY distance
		LIMIT $4
	`

	radiusMeters := radiusKm * 1000
	rows, err := r.db.GetConnection().QueryContext(ctx, query, lng, lat, radiusMeters, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search nearby places: %w", err)
	}
	defer rows.Close()

	var places []*models.Place
	for rows.Next() {
		place := &models.Place{}
		var distance float64
		err := rows.Scan(
			&place.ID, &place.Name, &place.Geometry, &place.Properties,
			&place.CreatedAt, &place.UpdatedAt, &place.DeletedAt, &distance,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan place: %w", err)
		}
		places = append(places, place)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate places: %w", err)
	}

	return places, nil
}

func (r *placeRepository) CreateRelation(ctx context.Context, relation *models.PlaceRelation) error {
	query := `
		INSERT INTO place_relations (id, from_place_id, to_place_id, relation_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		relation.ID, relation.FromPlaceID, relation.ToPlaceID, relation.RelationType,
		relation.CreatedAt, relation.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create place relation: %w", err)
	}

	return nil
}

func (r *placeRepository) GetRelations(ctx context.Context, placeID uuid.UUID) ([]*models.PlaceRelation, error) {
	query := `
		SELECT id, from_place_id, to_place_id, relation_type, created_at, updated_at, deleted_at
		FROM place_relations
		WHERE (from_place_id = $1 OR to_place_id = $1) AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, placeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get place relations: %w", err)
	}
	defer rows.Close()

	var relations []*models.PlaceRelation
	for rows.Next() {
		relation := &models.PlaceRelation{}
		err := rows.Scan(
			&relation.ID, &relation.FromPlaceID, &relation.ToPlaceID, &relation.RelationType,
			&relation.CreatedAt, &relation.UpdatedAt, &relation.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan place relation: %w", err)
		}
		relations = append(relations, relation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate place relations: %w", err)
	}

	return relations, nil
}

func (r *placeRepository) DeleteRelation(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE place_relations
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete place relation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("place relation not found or already deleted")
	}

	return nil
}
