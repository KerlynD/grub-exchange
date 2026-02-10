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
	notifRepo := repository.NewNotificationRepo(db)
	achieveRepo := repository.NewAchievementRepo(db)

	// Initialize services
	authService := services.NewAuthService(db, userRepo, balanceRepo, txnRepo)
	achieveSvc := services.NewAchievementService(achieveRepo, balanceRepo, portfolioRepo, userRepo)
	tradingService := services.NewTradingService(db, userRepo, balanceRepo, portfolioRepo, txnRepo, notifRepo, achieveSvc)
	portfolioService := services.NewPortfolioService(userRepo, balanceRepo, portfolioRepo, txnRepo)
	snapshotRepo := repository.NewMarketSnapshotRepo(db)
	marketService := services.NewMarketService(userRepo, balanceRepo, portfolioRepo, txnRepo, snapshotRepo)
	marketMaker := services.NewMarketMaker(db, userRepo, balanceRepo, portfolioRepo, txnRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	tradingHandler := handlers.NewTradingHandler(tradingService)
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)
	marketHandler := handlers.NewMarketHandler(marketService)
	profileHandler := handlers.NewProfileHandler(authService, userRepo)
	notifHandler := handlers.NewNotificationHandler(notifRepo)
	achieveHandler := handlers.NewAchievementHandler(achieveSvc)

	// Start background jobs
	go runScheduledJobs(marketService, achieveSvc, userRepo)
	go marketMaker.Run(60 * time.Second) // nudge prices every 60 seconds

	// Setup router
	router := api.SetupRouter(authHandler, tradingHandler, portfolioHandler, marketHandler, profileHandler, notifHandler, achieveHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Grub Exchange server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runScheduledJobs(marketService *services.MarketService, achieveSvc *services.AchievementService, userRepo *repository.UserRepo) {
	// Daily decay runs every 24h
	decayTicker := time.NewTicker(24 * time.Hour)
	defer decayTicker.Stop()

	// Dividends run daily (every 24 hours)
	dividendTicker := time.NewTicker(24 * time.Hour)
	defer dividendTicker.Stop()

	// Achievement check runs every hour (for Diamond Hands, Whale, etc.)
	achieveTicker := time.NewTicker(1 * time.Hour)
	defer achieveTicker.Stop()

	// Market snapshot every 5 minutes for the Grub Market chart
	snapshotTicker := time.NewTicker(5 * time.Minute)
	defer snapshotTicker.Stop()

	// Record initial snapshot on startup
	marketService.RecordMarketSnapshot()

	for {
		select {
		case <-decayTicker.C:
			log.Println("Running daily decay...")
			marketService.RunDailyDecay()
		case <-dividendTicker.C:
			log.Println("Running daily dividends...")
			marketService.RunDailyDividends()
		case <-achieveTicker.C:
			log.Println("Checking periodic achievements...")
			checkPeriodicAchievements(achieveSvc, userRepo)
		case <-snapshotTicker.C:
			marketService.RecordMarketSnapshot()
		}
	}
}

func checkPeriodicAchievements(achieveSvc *services.AchievementService, userRepo *repository.UserRepo) {
	users, err := userRepo.GetAll()
	if err != nil {
		log.Printf("Error checking achievements: %v", err)
		return
	}
	for _, u := range users {
		achieveSvc.CheckDiamondHands(u.ID)
		// Whale is also checked after trades, but check periodically too
		achieveSvc.CheckAfterTrade(u.ID)
	}
}
