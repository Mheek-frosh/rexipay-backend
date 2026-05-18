package controllers

import (
	"net/http"

	"github.com/Mheek-frosh/rexipaybackend/config"
	"github.com/Mheek-frosh/rexipaybackend/models"
	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var notifications []models.Notification

	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
