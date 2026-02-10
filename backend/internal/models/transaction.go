package models

import "time"

type Transaction struct {
	ID              int       `json:"id"`
	BuyerID         int       `json:"buyer_id"`
	StockUserID     int       `json:"stock_user_id"`
	TransactionType string    `json:"transaction_type"`
	NumShares       float64   `json:"num_shares"`
	PricePerShare   float64   `json:"price_per_share"`
	TotalGrub       float64   `json:"total_grub"`
	Timestamp       time.Time `json:"timestamp"`
}

type TransactionWithDetails struct {
	ID              int       `json:"id"`
	BuyerUsername   string    `json:"buyer_username"`
	StockTicker     string    `json:"stock_ticker"`
	TransactionType string    `json:"transaction_type"`
	NumShares       float64   `json:"num_shares"`
	PricePerShare   float64   `json:"price_per_share"`
	TotalGrub       float64   `json:"total_grub"`
	Timestamp       time.Time `json:"timestamp"`
}

type PriceHistory struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type TradeRequest struct {
	StockTicker string  `json:"stock_ticker" binding:"required"`
	NumShares   float64 `json:"num_shares"`
	GrubAmount  float64 `json:"grub_amount"`
}

type StockDetail struct {
	User              User                     `json:"user"`
	PriceHistory      []PriceHistory           `json:"price_history"`
	RecentTrades      []TransactionWithDetails `json:"recent_trades"`
	Change24h         float64                  `json:"change_24h"`
	Change24hPercent  float64                  `json:"change_24h_percent"`
	MarketCap         float64                  `json:"market_cap"`
	Volume24h         float64                  `json:"volume_24h"`
	AllTimeHigh       float64                  `json:"all_time_high"`
	AllTimeLow        float64                  `json:"all_time_low"`
}

type StockListItem struct {
	ID                int            `json:"id"`
	Username          string         `json:"username"`
	Ticker            string         `json:"ticker"`
	CurrentSharePrice float64        `json:"current_share_price"`
	Change24hPercent  float64        `json:"change_24h_percent"`
	SparklineData     []float64      `json:"sparkline_data"`
}

type LeaderboardData struct {
	MostValuable    []LeaderboardEntry `json:"most_valuable"`
	BiggestGainers  []LeaderboardEntry `json:"biggest_gainers"`
	BiggestLosers   []LeaderboardEntry `json:"biggest_losers"`
	RichestTraders  []LeaderboardEntry `json:"richest_traders"`
	BestPerformance []LeaderboardEntry `json:"best_performance"`
}

type LeaderboardEntry struct {
	Rank     int     `json:"rank"`
	UserID   int     `json:"user_id"`
	Username string  `json:"username"`
	Ticker   string  `json:"ticker"`
	Value    float64 `json:"value"`
	Change   float64 `json:"change,omitempty"`
}

type MarketOverview struct {
	TotalMarketCap  float64 `json:"total_market_cap"`
	TotalGrub       float64 `json:"total_grub"`
	TotalInvested   float64 `json:"total_invested"`
	TotalCash       float64 `json:"total_cash"`
	InvestedPercent float64 `json:"invested_percent"`
	TotalStocks     int     `json:"total_stocks"`
}
