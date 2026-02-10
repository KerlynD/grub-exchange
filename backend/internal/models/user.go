package models

import "time"

type User struct {
	ID                int        `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	Ticker            string     `json:"ticker"`
	Bio               string     `json:"bio"`
	CurrentSharePrice float64    `json:"current_share_price"`
	SharesOutstanding int        `json:"shares_outstanding"`
	LastLogin         *time.Time `json:"last_login,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=2,max=15"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID                int        `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	Ticker            string     `json:"ticker"`
	Bio               string     `json:"bio"`
	CurrentSharePrice float64    `json:"current_share_price"`
	SharesOutstanding int        `json:"shares_outstanding"`
	GrubBalance       float64    `json:"grub_balance"`
	LastLogin         *time.Time `json:"last_login,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type UpdateProfileRequest struct {
	Bio string `json:"bio" binding:"max=500"`
}

type PortfolioSnapshot struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	TotalValue  float64   `json:"total_value"`
	GrubBalance float64   `json:"grub_balance"`
	Timestamp   time.Time `json:"timestamp"`
}
