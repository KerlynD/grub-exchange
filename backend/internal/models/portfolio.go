package models

type Portfolio struct {
	ID               int     `json:"id"`
	OwnerID          int     `json:"owner_id"`
	StockUserID      int     `json:"stock_user_id"`
	NumShares        float64 `json:"num_shares"`
	AvgPurchasePrice float64 `json:"avg_purchase_price"`
}

type PortfolioHolding struct {
	Ticker            string  `json:"ticker"`
	Username          string  `json:"username"`
	StockUserID       int     `json:"stock_user_id"`
	NumShares         float64 `json:"num_shares"`
	AvgPurchasePrice  float64 `json:"avg_purchase_price"`
	CurrentPrice      float64 `json:"current_price"`
	TotalValue        float64 `json:"total_value"`
	ProfitLoss        float64 `json:"profit_loss"`
	ProfitLossPercent float64 `json:"profit_loss_percent"`
}

type PortfolioResponse struct {
	GrubBalance     float64            `json:"grub_balance"`
	TotalValue      float64            `json:"total_value"`
	TotalPL         float64            `json:"total_pl"`
	TotalPLPercent  float64            `json:"total_pl_percent"`
	Holdings        []PortfolioHolding `json:"holdings"`
	CanClaimDaily   bool               `json:"can_claim_daily"`
}
