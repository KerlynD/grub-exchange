package api

import (
	"grub-exchange/internal/api/handlers"
	"grub-exchange/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	tradingHandler *handlers.TradingHandler,
	portfolioHandler *handlers.PortfolioHandler,
	marketHandler *handlers.MarketHandler,
	profileHandler *handlers.ProfileHandler,
	notifHandler *handlers.NotificationHandler,
	achieveHandler *handlers.AchievementHandler,
	postHandler *handlers.PostHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/auth/me", authHandler.GetMe)

			// Trading
			protected.POST("/trade/buy", tradingHandler.Buy)
			protected.POST("/trade/sell", tradingHandler.Sell)

			// Portfolio
			protected.GET("/portfolio", portfolioHandler.GetPortfolio)
			protected.GET("/portfolio/history", portfolioHandler.GetHistory)
			protected.POST("/portfolio/claim-daily", portfolioHandler.ClaimDaily)
			protected.GET("/portfolio/graph", profileHandler.GetPortfolioGraph)

			// Profile
			protected.GET("/profile", profileHandler.GetProfile)
			protected.PUT("/profile", profileHandler.UpdateProfile)

			// Market
			protected.GET("/market/overview", marketHandler.GetMarketOverview)
			protected.GET("/stocks", marketHandler.GetStocks)
			protected.GET("/stocks/:ticker", marketHandler.GetStockDetail)
			protected.GET("/leaderboard", marketHandler.GetLeaderboard)
			protected.GET("/transactions", marketHandler.GetRecentTransactions)

			// Notifications
			protected.GET("/notifications", notifHandler.GetNotifications)
			protected.POST("/notifications/read", notifHandler.MarkRead)

			// Achievements
			protected.GET("/achievements", achieveHandler.GetMyAchievements)

			// News / Posts
			protected.GET("/posts/recent", postHandler.GetRecentPosts)
			protected.GET("/stocks/:ticker/posts", postHandler.GetPosts)
			protected.POST("/stocks/:ticker/posts", postHandler.CreatePost)
			protected.POST("/posts/:id/vote", postHandler.VotePost)
		}
	}

	return r
}
