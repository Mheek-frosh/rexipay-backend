package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mheek-frosh/rexipaybackend/config"
	"github.com/Mheek-frosh/rexipaybackend/models"
	"github.com/gin-gonic/gin"
)

type DepositInput struct {
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	WalletType string  `json:"wallet_type" binding:"required,oneof=fiat crypto"`
}

func Deposit(c *gin.Context) {
	var input DepositInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	// 1. Find Wallet
	var wallet models.Wallet
	if err := config.DB.Where("user_id = ? AND type = ?", userID, input.WalletType).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// 2. Begin Transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. Update Balance
	wallet.Balance += input.Amount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
		return
	}

	// 4. Create Transaction Record
	transaction := models.Transaction{
		SenderWalletID:   nil, // System/External deposit (no sender wallet)
		ReceiverWalletID: &wallet.ID,
		Amount:           input.Amount,
		TransactionType:  "deposit",
		Status:           "success",
		Reference:        fmt.Sprintf("DEP-%d", time.Now().UnixNano()),
		CreatedAt:        time.Now(),
	}

	// Adapting to existing GORM model if it differs.
	// The User SQL has `receiver_wallet_id`. The existing Go model might be different.
	// I will check the Transaction model in a moment, but for now logic is standard.

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	// 5. Create Notification
	notification := models.Notification{
		UserID:  uint(userID.(float64)),
		Title:   "Deposit Successful",
		Message: fmt.Sprintf("Your wallet has been credited with %f", input.Amount),
	}
	tx.Create(&notification)

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Deposit successful",
		"balance":     wallet.Balance,
		"transaction": transaction,
	})
}
