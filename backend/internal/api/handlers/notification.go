package handlers

import (
	"grub-exchange/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifRepo *repository.NotificationRepo
}

func NewNotificationHandler(notifRepo *repository.NotificationRepo) *NotificationHandler {
	return &NotificationHandler{notifRepo: notifRepo}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	notifs, err := h.notifRepo.GetByUser(userID, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unread, _ := h.notifRepo.GetUnreadCount(userID)

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifs,
		"unread_count":  unread,
	})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	if err := h.notifRepo.MarkAllRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}
