package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/models"
)

type notificationRepository struct {
	db *database.DB
}

func NewNotificationRepository(db *database.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	if notification.Data == nil {
		notification.Data = map[string]string{}
	}

	dataJSON, err := json.Marshal(notification.Data)
	if err != nil {
		return fmt.Errorf("marshal notification data: %w", err)
	}

	query := `
		INSERT INTO notifications (id, user_id, actor_id, post_id, comment_id, type, data, read_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = r.db.GetConnection().ExecContext(ctx, query,
		notification.ID,
		notification.UserID,
		notification.ActorID,
		notification.PostID,
		notification.CommentID,
		notification.Type,
		dataJSON,
		notification.ReadAt,
		notification.CreatedAt,
		notification.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	return nil
}

func (r *notificationRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	query := `
		SELECT id, user_id, actor_id, post_id, comment_id, type, data, read_at, created_at, updated_at, deleted_at
		FROM notifications
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var notif models.Notification
		var dataJSON []byte
		if err := rows.Scan(
			&notif.ID,
			&notif.UserID,
			&notif.ActorID,
			&notif.PostID,
			&notif.CommentID,
			&notif.Type,
			&dataJSON,
			&notif.ReadAt,
			&notif.CreatedAt,
			&notif.UpdatedAt,
			&notif.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}

		if len(dataJSON) > 0 {
			if err := json.Unmarshal(dataJSON, &notif.Data); err != nil {
				return nil, fmt.Errorf("unmarshal notification data: %w", err)
			}
		}

		notifications = append(notifications, &notif)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notifications: %w", err)
	}

	return notifications, nil
}

func (r *notificationRepository) ClearByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET deleted_at = NOW()
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	if _, err := r.db.GetConnection().ExecContext(ctx, query, userID); err != nil {
		return fmt.Errorf("clear notifications: %w", err)
	}

	return nil
}

func (r *notificationRepository) SoftDeleteByReference(ctx context.Context, userID, actorID uuid.UUID, notifType models.NotificationType, postID *uuid.UUID, commentID *uuid.UUID) error {
	query := `
		UPDATE notifications
		SET deleted_at = NOW()
		WHERE user_id = $1 AND actor_id = $2 AND type = $3
			AND (post_id = $4 OR ($4 IS NULL AND post_id IS NULL))
			AND (comment_id = $5 OR ($5 IS NULL AND comment_id IS NULL))
			AND deleted_at IS NULL
	`

	if _, err := r.db.GetConnection().ExecContext(ctx, query, userID, actorID, notifType, postID, commentID); err != nil {
		return fmt.Errorf("soft delete notification: %w", err)
	}

	return nil
}
