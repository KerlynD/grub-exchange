package handlers

import (
	"grub-exchange/internal/models"
	"grub-exchange/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TradingHandler struct {
	tradingService *services.TradingService
}

func NewTradingHandler(tradingService *services.TradingService) *TradingHandler {
	return &TradingHandler{tradingService: tradingService}
}

func (h *TradingHandler) Buy(c *gin.Context) {
	var req models.TradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	if req.NumShares <= 0 && req.GrubAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "specify either num_shares or grub_amount"})
		return
	}

	userID, _ := c.Get("userID")

	txn, err := h.tradingService.ExecuteBuy(userID.(int), req.StockTicker, req.NumShares, req.GrubAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "buy order executed",
		"transaction": txn,
	})
}

func (h *TradingHandler) Sell(c *gin.Context) {
	var req models.TradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	if req.NumShares <= 0 && req.GrubAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "specify either num_shares or grub_amount"})
		return
	}

	userID, _ := c.Get("userID")

	txn, err := h.tradingService.ExecuteSell(userID.(int), req.StockTicker, req.NumShares, req.GrubAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "sell order executed",
		"transaction": txn,
	})
}
