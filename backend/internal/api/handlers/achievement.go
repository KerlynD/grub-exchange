package handlers

import (
	"grub-exchange/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AchievementHandler struct {
	achievementService *services.AchievementService
}

func NewAchievementHandler(achievementService *services.AchievementService) *AchievementHandler {
	return &AchievementHandler{achievementService: achievementService}
}

func (h *AchievementHandler) GetMyAchievements(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	earned, err := h.achievementService.GetUserAchievements(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	all, err := h.achievementService.GetAllAchievements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"earned": earned,
		"all":    all,
	})
}
