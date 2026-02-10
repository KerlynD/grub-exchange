package services

import (
	"database/sql"
	"errors"
	"fmt"
	"grub-exchange/internal/models"
	"grub-exchange/internal/repository"
	"math"
)

type TradingService struct {
	db            *sql.DB
	userRepo      *repository.UserRepo
	balanceRepo   *repository.BalanceRepo
	portfolioRepo *repository.PortfolioRepo
	txnRepo       *repository.TransactionRepo
	notifRepo     *repository.NotificationRepo
	achieveSvc    *AchievementService
}

func NewTradingService(
	db *sql.DB,
	userRepo *repository.UserRepo,
	balanceRepo *repository.BalanceRepo,
	portfolioRepo *repository.PortfolioRepo,
	txnRepo *repository.TransactionRepo,
	notifRepo *repository.NotificationRepo,
	achieveSvc *AchievementService,
) *TradingService {
	return &TradingService{
		db:            db,
		userRepo:      userRepo,
		balanceRepo:   balanceRepo,
		portfolioRepo: portfolioRepo,
		txnRepo:       txnRepo,
		notifRepo:     notifRepo,
		achieveSvc:    achieveSvc,
	}
}

// ResolveShares converts either a share count or grub amount into a final share count.
// When grubAmount is specified, shares are calculated using the EXECUTION price (not spot)
// so the user actually gets the number of shares their grub can buy after slippage.
func ResolveShares(numShares, grubAmount, executionPrice float64) (float64, error) {
	if grubAmount > 0 {
		if executionPrice <= 0 {
			return 0, errors.New("invalid stock price")
		}
		shares := grubAmount / executionPrice
		// Round to 4 decimal places
		shares = math.Round(shares*10000) / 10000
		if shares <= 0 {
			return 0, errors.New("amount too small to buy any shares")
		}
		return shares, nil
	}
	if numShares <= 0 {
		return 0, errors.New("specify either num_shares or grub_amount")
	}
	// Round to 4 decimal places
	return math.Round(numShares*10000) / 10000, nil
}

func (s *TradingService) ExecuteBuy(buyerID int, stockTicker string, numShares float64, grubAmount float64) (*models.TransactionWithDetails, error) {
	stockUser, err := s.userRepo.GetByTicker(stockTicker)
	if err != nil {
		return nil, errors.New("stock not found")
	}

	if stockUser.ID == buyerID {
		return nil, errors.New("cannot buy your own stock")
	}

	// Calculate execution price FIRST so both share-based and grub-based orders
	// use the same fair price that includes price impact (slippage).
	// For grub-based orders, we need to estimate shares first to get exec price,
	// then re-resolve with the exec price.
	var finalShares float64
	var execPrice float64
	var newPrice float64

	if grubAmount > 0 {
		// Estimate shares at spot price, then get exec price, then re-resolve
		estShares := grubAmount / stockUser.CurrentSharePrice
		estShares = math.Round(estShares*10000) / 10000
		newPrice, execPrice = CalculateTradeExecution(stockUser.CurrentSharePrice, estShares, float64(stockUser.SharesOutstanding))
		finalShares, err = ResolveShares(0, grubAmount, execPrice)
		if err != nil {
			return nil, err
		}
		// Recalculate with actual shares
		newPrice, execPrice = CalculateTradeExecution(stockUser.CurrentSharePrice, finalShares, float64(stockUser.SharesOutstanding))
	} else {
		finalShares, err = ResolveShares(numShares, 0, stockUser.CurrentSharePrice)
		if err != nil {
			return nil, err
		}
		newPrice, execPrice = CalculateTradeExecution(stockUser.CurrentSharePrice, finalShares, float64(stockUser.SharesOutstanding))
	}

	balance, err := s.balanceRepo.GetByUserID(buyerID)
	if err != nil {
		return nil, errors.New("balance not found")
	}

	// Cost is based on execution price (includes slippage), not spot price
	totalCost := finalShares * execPrice
	totalCost = math.Round(totalCost*100) / 100
	if balance.GrubBalance < totalCost {
		return nil, errors.New("insufficient Grub balance")
	}

	// Read existing holding BEFORE starting transaction
	existing, err := s.portfolioRepo.GetHolding(buyerID, stockUser.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var newAvgPrice float64
	var newSharesTotal float64
	if existing != nil {
		newSharesTotal = existing.NumShares + finalShares
		newAvgPrice = ((existing.AvgPurchasePrice * existing.NumShares) + (execPrice * finalShares)) / newSharesTotal
	} else {
		newSharesTotal = finalShares
		newAvgPrice = execPrice
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.balanceRepo.UpdateBalance(tx, buyerID, -totalCost); err != nil {
		return nil, err
	}

	if err := s.portfolioRepo.UpsertHolding(tx, buyerID, stockUser.ID, newSharesTotal, newAvgPrice); err != nil {
		return nil, err
	}

	if err := s.userRepo.UpdateSharePrice(tx, stockUser.ID, newPrice); err != nil {
		return nil, err
	}

	// 2% appreciation to stock owner (based on execution cost)
	appreciation := totalCost * 0.02
	if err := s.balanceRepo.UpdateBalance(tx, stockUser.ID, appreciation); err != nil {
		return nil, err
	}

	txn := &models.Transaction{
		BuyerID:         buyerID,
		StockUserID:     stockUser.ID,
		TransactionType: "BUY",
		NumShares:       finalShares,
		PricePerShare:   execPrice,
		TotalGrub:       totalCost,
	}
	if err := s.txnRepo.Create(tx, txn); err != nil {
		return nil, err
	}

	if err := s.txnRepo.RecordPriceHistory(tx, stockUser.ID, newPrice); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	buyer, _ := s.userRepo.GetByID(buyerID)
	buyerUsername := ""
	if buyer != nil {
		buyerUsername = buyer.Username
	}

	// Notify the stock owner that someone bought their stock
	if s.notifRepo != nil && stockUser.ID != buyerID {
		msg := fmt.Sprintf("%s just bought %.2f shares of you!", buyerUsername, finalShares)
		_ = s.notifRepo.Create(stockUser.ID, "trade_buy", msg, buyerUsername, stockUser.Ticker, finalShares)
	}

	// Check achievements for the buyer
	if s.achieveSvc != nil {
		s.achieveSvc.CheckAfterTrade(buyerID)
	}

	return &models.TransactionWithDetails{
		BuyerUsername:    buyerUsername,
		StockTicker:     stockUser.Ticker,
		TransactionType: "BUY",
		NumShares:       finalShares,
		PricePerShare:   execPrice,
		TotalGrub:       totalCost,
	}, nil
}

func (s *TradingService) ExecuteSell(sellerID int, stockTicker string, numShares float64, grubAmount float64) (*models.TransactionWithDetails, error) {
	stockUser, err := s.userRepo.GetByTicker(stockTicker)
	if err != nil {
		return nil, errors.New("stock not found")
	}

	// Calculate execution price for sells too — seller eats downward slippage
	var finalShares float64
	var execPrice float64
	var newPrice float64

	if grubAmount > 0 {
		// Selling by grub amount: estimate shares first
		estShares := grubAmount / stockUser.CurrentSharePrice
		estShares = math.Round(estShares*10000) / 10000
		newPrice, execPrice = CalculateTradeExecution(stockUser.CurrentSharePrice, -estShares, float64(stockUser.SharesOutstanding))
		finalShares, err = ResolveShares(0, grubAmount, execPrice)
		if err != nil {
			return nil, err
		}
		newPrice, execPrice = CalculateTradeExecution(stockUser.CurrentSharePrice, -finalShares, float64(stockUser.SharesOutstanding))
	} else {
		finalShares, err = ResolveShares(numShares, 0, stockUser.CurrentSharePrice)
		if err != nil {
			return nil, err
		}
		newPrice, execPrice = CalculateTradeExecution(stockUser.CurrentSharePrice, -finalShares, float64(stockUser.SharesOutstanding))
	}

	holding, err := s.portfolioRepo.GetHolding(sellerID, stockUser.ID)
	if err != nil {
		return nil, errors.New("you don't own any shares of this stock")
	}

	if holding.NumShares < finalShares {
		return nil, errors.New("insufficient shares to sell")
	}

	// Proceeds based on execution price (includes downward slippage)
	totalProceeds := finalShares * execPrice
	totalProceeds = math.Round(totalProceeds*100) / 100

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.balanceRepo.UpdateBalance(tx, sellerID, totalProceeds); err != nil {
		return nil, err
	}

	remainingShares := holding.NumShares - finalShares
	if remainingShares <= 0.0001 {
		if err := s.portfolioRepo.DeleteHolding(tx, sellerID, stockUser.ID); err != nil {
			return nil, err
		}
	} else {
		if err := s.portfolioRepo.UpsertHolding(tx, sellerID, stockUser.ID, remainingShares, holding.AvgPurchasePrice); err != nil {
			return nil, err
		}
	}

	if err := s.userRepo.UpdateSharePrice(tx, stockUser.ID, newPrice); err != nil {
		return nil, err
	}

	txn := &models.Transaction{
		BuyerID:         sellerID,
		StockUserID:     stockUser.ID,
		TransactionType: "SELL",
		NumShares:       finalShares,
		PricePerShare:   execPrice,
		TotalGrub:       totalProceeds,
	}
	if err := s.txnRepo.Create(tx, txn); err != nil {
		return nil, err
	}

	if err := s.txnRepo.RecordPriceHistory(tx, stockUser.ID, newPrice); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	seller, _ := s.userRepo.GetByID(sellerID)
	sellerUsername := ""
	if seller != nil {
		sellerUsername = seller.Username
	}

	// Notify the stock owner that someone sold their stock
	if s.notifRepo != nil && stockUser.ID != sellerID {
		msg := fmt.Sprintf("%s just sold %.2f shares of you", sellerUsername, finalShares)
		_ = s.notifRepo.Create(stockUser.ID, "trade_sell", msg, sellerUsername, stockUser.Ticker, finalShares)
	}

	// Check achievements for the seller
	if s.achieveSvc != nil {
		s.achieveSvc.CheckAfterTrade(sellerID)
	}

	return &models.TransactionWithDetails{
		BuyerUsername:    sellerUsername,
		StockTicker:     stockUser.Ticker,
		TransactionType: "SELL",
		NumShares:       finalShares,
		PricePerShare:   execPrice,
		TotalGrub:       totalProceeds,
	}, nil
}
