package main

import (
	"grub-exchange/internal/api"
	"grub-exchange/internal/api/handlers"
	"grub-exchange/internal/database"
	"grub-exchange/internal/repository"
	"grub-exchange/internal/services"
	"log"
	"os"
	"time"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepo(db)
	balanceRepo := repository.NewBalanceRepo(db)
	portfolioRepo := repository.NewPortfolioRepo(db)
	txnRepo := repository.NewTransactionRepo(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, balanceRepo, txnRepo)
	tradingService := services.NewTradingService(db, userRepo, balanceRepo, portfolioRepo, txnRepo)
	portfolioService := services.NewPortfolioService(userRepo, balanceRepo, portfolioRepo, txnRepo)
	marketService := services.NewMarketService(userRepo, balanceRepo, portfolioRepo, txnRepo)
	marketMaker := services.NewMarketMaker(userRepo, balanceRepo, portfolioRepo, txnRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	tradingHandler := handlers.NewTradingHandler(tradingService)
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)
	marketHandler := handlers.NewMarketHandler(marketService)
	profileHandler := handlers.NewProfileHandler(authService, userRepo)

	// Start background jobs
	go runScheduledJobs(marketService)
	go marketMaker.Run(30 * time.Second) // nudge prices every 30 seconds

	// Setup router
	router := api.SetupRouter(authHandler, tradingHandler, portfolioHandler, marketHandler, profileHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Grub Exchange server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runScheduledJobs(marketService *services.MarketService) {
	// Daily decay runs every 24h
	decayTicker := time.NewTicker(24 * time.Hour)
	defer decayTicker.Stop()

	// Dividends run biweekly (every 14 days)
	dividendTicker := time.NewTicker(14 * 24 * time.Hour)
	defer dividendTicker.Stop()

	for {
		select {
		case <-decayTicker.C:
			log.Println("Running daily decay...")
			marketService.RunDailyDecay()
		case <-dividendTicker.C:
			log.Println("Running biweekly dividends...")
			marketService.RunDailyDividends()
		}
	}
}
