package repository

import (
	"database/sql"
	"grub-exchange/internal/models"
	"time"
)

type NotificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(userID int, notifType, message, actorUsername, stockTicker string, numShares float64) error {
	_, err := r.db.Exec(
		`INSERT INTO notifications (user_id, type, message, actor_username, stock_ticker, num_shares, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, notifType, message, actorUsername, stockTicker, numShares, time.Now(),
	)
	return err
}

func (r *NotificationRepo) GetByUser(userID int, limit int) ([]models.Notification, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, type, message, actor_username, stock_ticker, num_shares, read, created_at
		 FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Message, &n.ActorUsername, &n.StockTicker, &n.NumShares, &n.Read, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (r *NotificationRepo) MarkAllRead(userID int) error {
	_, err := r.db.Exec(
		`UPDATE notifications SET read = true WHERE user_id = $1 AND read = false`,
		userID,
	)
	return err
}

func (r *NotificationRepo) GetUnreadCount(userID int) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false`,
		userID,
	).Scan(&count)
	return count, err
}
