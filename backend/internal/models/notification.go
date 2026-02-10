package models

import "time"

type Notification struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Type          string    `json:"type"`
	Message       string    `json:"message"`
	ActorUsername string    `json:"actor_username"`
	StockTicker   string    `json:"stock_ticker"`
	NumShares     float64   `json:"num_shares"`
	Read          bool      `json:"read"`
	CreatedAt     time.Time `json:"created_at"`
}

type Achievement struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type UserAchievement struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	AchievementID string    `json:"achievement_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Icon          string    `json:"icon"`
	EarnedAt      time.Time `json:"earned_at"`
}
