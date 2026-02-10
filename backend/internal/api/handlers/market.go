package handlers

import (
	"grub-exchange/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MarketHandler struct {
	marketService *services.MarketService
}

func NewMarketHandler(marketService *services.MarketService) *MarketHandler {
	return &MarketHandler{marketService: marketService}
}

func (h *MarketHandler) GetStocks(c *gin.Context) {
	stocks, err := h.marketService.GetAllStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

func (h *MarketHandler) GetStockDetail(c *gin.Context) {
	ticker := c.Param("ticker")
	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticker is required"})
		return
	}

	detail, err := h.marketService.GetStockDetail(ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
		return
	}

	c.JSON(http.StatusOK, detail)
}

func (h *MarketHandler) GetLeaderboard(c *gin.Context) {
	data, err := h.marketService.GetLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *MarketHandler) GetRecentTransactions(c *gin.Context) {
	txns, err := h.marketService.GetRecentTransactions(20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": txns})
}
