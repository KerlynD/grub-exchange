package handlers

import (
	"grub-exchange/internal/models"
	"grub-exchange/internal/repository"
	"grub-exchange/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	authService *services.AuthService
	userRepo    *repository.UserRepo
}

func NewProfileHandler(authService *services.AuthService, userRepo *repository.UserRepo) *ProfileHandler {
	return &ProfileHandler{authService: authService, userRepo: userRepo}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	user, err := h.authService.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	if err := h.userRepo.UpdateBio(userID, req.Bio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	user, err := h.authService.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *ProfileHandler) GetPortfolioGraph(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	since := time.Now().Add(-90 * 24 * time.Hour)
	snapshots, err := h.userRepo.GetPortfolioSnapshots(userID, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"snapshots": snapshots})
}
