package controllers

import (
	"net/http"
	"time"

	"github.com/Mheek-frosh/rexipaybackend/config"
	"github.com/Mheek-frosh/rexipaybackend/models"
	"github.com/gin-gonic/gin"
)

type AddBeneficiaryInput struct {
	BankCode      string `json:"bank_code" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required"`
	AccountName   string `json:"account_name" binding:"required"`
}

func AddBeneficiary(c *gin.Context) {
	var input AddBeneficiaryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uid := uint(userID.(float64))

	beneficiary := models.Beneficiary{
		UserID:        uid,
		BankCode:      input.BankCode,
		AccountNumber: input.AccountNumber,
		AccountName:   input.AccountName,
		CreatedAt:     time.Now(),
	}

	if err := config.DB.Create(&beneficiary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add beneficiary"})
		return
	}

	c.JSON(http.StatusOK, beneficiary)
}

func GetBeneficiaries(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var beneficiaries []models.Beneficiary

	if err := config.DB.Where("user_id = ?", userID).Find(&beneficiaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch beneficiaries"})
		return
	}

	c.JSON(http.StatusOK, beneficiaries)
}

func DeleteBeneficiary(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Beneficiary{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete beneficiary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beneficiary deleted successfully"})
}
