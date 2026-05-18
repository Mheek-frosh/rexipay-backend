package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mheek-frosh/rexipaybackend/config"
	"github.com/Mheek-frosh/rexipaybackend/models"
	"github.com/Mheek-frosh/rexipaybackend/services"
	"github.com/Mheek-frosh/rexipaybackend/utils"
	"github.com/gin-gonic/gin"
)

type CreateWalletInput struct {
	Type string `json:"type" binding:"required,oneof=fiat crypto"`
}

type LookupInput struct {
	AccountNumber string `json:"account_number" binding:"required"`
	BankCode      string `json:"bank_code" binding:"required"`
}

type TransferInput struct {
	RecipientBank string  `json:"recipient_bank"`
	RecipientAcc  string  `json:"recipient_account" binding:"required"`
	RecipientName string  `json:"recipient_name" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	WalletType    string  `json:"wallet_type" binding:"required,oneof=fiat crypto"`
}

func CreateWallet(c *gin.Context) {
	var input CreateWalletInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	currency := "NGN"
	if input.Type == "crypto" {
		currency = "USDT"
	}

	wallet := models.Wallet{
		UserID:        uint(userID.(float64)),
		Type:          input.Type,
		Currency:      currency,
		Balance:       0,
		AccountNumber: utils.GenerateAccountNumber(),
		BankName:      "Rexi Bank",
	}

	if err := config.DB.Create(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func GetUserWallets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var wallets []models.Wallet
	if err := config.DB.Where("user_id = ?", userID).Find(&wallets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wallets"})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

func GetTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 1. Get User's Wallet IDs
	var walletIDs []uint
	if err := config.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("id", &walletIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wallets"})
		return
	}

	if len(walletIDs) == 0 {
		c.JSON(http.StatusOK, []models.Transaction{})
		return
	}

	var transactions []models.Transaction
	// Fetch transactions where the user's wallet is sender OR receiver
	if err := config.DB.Where("sender_wallet_id IN ? OR receiver_wallet_id IN ?", walletIDs, walletIDs).Order("created_at desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func LookupAccount(c *gin.Context) {
	var input LookupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call Paystack Service to resolve account
	accountName, err := services.ResolveAccountNumber(input.AccountNumber, input.BankCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not verify account: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_number": input.AccountNumber,
		"bank_code":      input.BankCode,
		"account_name":   accountName,
	})
}

func TransferFunds(c *gin.Context) {
	var input TransferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	// 1. Find the correct wallet
	var wallet models.Wallet
	if err := config.DB.Where("user_id = ? AND type = ?", userID, input.WalletType).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// 2. Check balance
	if wallet.Balance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// 3. Perform transfer (DB Transaction)
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Deduct from wallet
	wallet.Balance -= input.Amount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
		return
	}

	// Create transaction record
	transaction := models.Transaction{
		SenderWalletID:  &wallet.ID,
		RecipientBank:   input.RecipientBank,
		RecipientAcc:    input.RecipientAcc,
		RecipientName:   input.RecipientName,
		Amount:          input.Amount,
		TransactionType: "transfer",
		Status:          "success",
		Reference:       fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		CreatedAt:       time.Now(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	// Create Notification
	notification := models.Notification{
		UserID:  uint(userID.(float64)),
		Title:   "Debit Alert",
		Message: fmt.Sprintf("You have successfully transferred %f to %s", input.Amount, input.RecipientName),
	}
	tx.Create(&notification)

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":        "Transfer successful",
		"transaction_id": transaction.ID,
		"amount":         transaction.Amount,
		"reference":      transaction.Reference,
		"time":           transaction.CreatedAt,
	})
}
