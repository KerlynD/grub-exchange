package repository

import (
	"database/sql"
	"grub-exchange/internal/models"
	"time"
)

type AchievementRepo struct {
	db *sql.DB
}

func NewAchievementRepo(db *sql.DB) *AchievementRepo {
	return &AchievementRepo{db: db}
}

func (r *AchievementRepo) Award(userID int, achievementID string) error {
	_, err := r.db.Exec(
		`INSERT INTO user_achievements (user_id, achievement_id, earned_at) VALUES ($1, $2, $3) ON CONFLICT (user_id, achievement_id) DO NOTHING`,
		userID, achievementID, time.Now(),
	)
	return err
}

func (r *AchievementRepo) HasAchievement(userID int, achievementID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM user_achievements WHERE user_id = $1 AND achievement_id = $2`,
		userID, achievementID,
	).Scan(&count)
	return count > 0, err
}

func (r *AchievementRepo) GetByUser(userID int) ([]models.UserAchievement, error) {
	rows, err := r.db.Query(
		`SELECT ua.id, ua.user_id, ua.achievement_id, a.name, a.description, a.icon, ua.earned_at
		 FROM user_achievements ua
		 JOIN achievements a ON ua.achievement_id = a.id
		 WHERE ua.user_id = $1 ORDER BY ua.earned_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []models.UserAchievement
	for rows.Next() {
		var a models.UserAchievement
		if err := rows.Scan(&a.ID, &a.UserID, &a.AchievementID, &a.Name, &a.Description, &a.Icon, &a.EarnedAt); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

func (r *AchievementRepo) GetAll() ([]models.Achievement, error) {
	rows, err := r.db.Query(`SELECT id, name, description, icon FROM achievements ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []models.Achievement
	for rows.Next() {
		var a models.Achievement
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Icon); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

// GetTradeCountToday returns the number of trades a user made today
func (r *AchievementRepo) GetTradeCountToday(userID int) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM transactions WHERE buyer_id = $1 AND timestamp >= CURRENT_DATE`,
		userID,
	).Scan(&count)
	return count, err
}

// GetOldestHoldingDays returns the age in days of the user's oldest holding
func (r *AchievementRepo) GetOldestHoldingDays(userID int) (int, error) {
	var days sql.NullInt64
	err := r.db.QueryRow(
		`SELECT EXTRACT(DAY FROM NOW() - MIN(t.timestamp))::INTEGER
		 FROM portfolios p
		 JOIN transactions t ON t.buyer_id = p.owner_id AND t.stock_user_id = p.stock_user_id AND t.transaction_type = 'BUY'
		 WHERE p.owner_id = $1 AND p.num_shares > 0`,
		userID,
	).Scan(&days)
	if err != nil || !days.Valid {
		return 0, err
	}
	return int(days.Int64), nil
}
