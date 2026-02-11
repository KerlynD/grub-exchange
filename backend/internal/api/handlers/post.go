package handlers

import (
	"grub-exchange/internal/models"
	"grub-exchange/internal/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postRepo *repository.PostRepo
	userRepo *repository.UserRepo
}

func NewPostHandler(postRepo *repository.PostRepo, userRepo *repository.UserRepo) *PostHandler {
	return &PostHandler{postRepo: postRepo, userRepo: userRepo}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	ticker := c.Param("ticker")
	stockUser, err := h.userRepo.GetByTicker(ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" || len(content) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content must be 1-500 characters"})
		return
	}

	post, err := h.postRepo.Create(userID, stockUser.ID, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Fill in the author username
	author, _ := h.userRepo.GetByID(userID)
	if author != nil {
		post.AuthorUsername = author.Username
	}
	post.StockTicker = stockUser.Ticker

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	ticker := c.Param("ticker")
	stockUser, err := h.userRepo.GetByTicker(ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	postsResult, err := h.postRepo.GetByStock(stockUser.ID, userID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	posts := postsResult
	if posts == nil {
		posts = []models.StockPost{}
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *PostHandler) VotePost(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req struct {
		VoteType int `json:"vote_type"` // 1 = like, -1 = dislike
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.VoteType != 1 && req.VoteType != -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vote_type must be 1 (like) or -1 (dislike)"})
		return
	}

	if err := h.postRepo.Vote(postID, userID, req.VoteType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded"})
}

func (h *PostHandler) GetRecentPosts(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	postsResult, err := h.postRepo.GetRecent(userID, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	posts := postsResult
	if posts == nil {
		posts = []models.StockPost{}
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
