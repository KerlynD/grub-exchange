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
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	portfolio, err := h.portfolioService.GetUserPortfolio(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

func (h *PortfolioHandler) ClaimDaily(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	newBalance, err := h.portfolioService.ClaimDailyBonus(userID)
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
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	history, err := h.portfolioService.GetTransactionHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": history})
}
