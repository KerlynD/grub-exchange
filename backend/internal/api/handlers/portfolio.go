package handlers

import (
	"grub-exchange/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	portfolioService *services.PortfolioService
}

func NewPortfolioHandler(portfolioService *services.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{portfolioService: portfolioService}
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	userID, _ := c.Get("userID")

	portfolio, err := h.portfolioService.GetUserPortfolio(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

func (h *PortfolioHandler) ClaimDaily(c *gin.Context) {
	userID, _ := c.Get("userID")

	newBalance, err := h.portfolioService.ClaimDailyBonus(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "daily bonus claimed!",
		"new_balance": newBalance,
	})
}

func (h *PortfolioHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

	history, err := h.portfolioService.GetTransactionHistory(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": history})
}
