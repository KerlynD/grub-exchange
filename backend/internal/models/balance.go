package models

import "time"

type Balance struct {
	UserID         int        `json:"user_id"`
	GrubBalance    float64    `json:"grub_balance"`
	LastDailyClaim *time.Time `json:"last_daily_claim,omitempty"`
}
