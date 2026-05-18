package controllers

import (
	"net/http"

	"github.com/Mheek-frosh/rexipaybackend/config"
	"github.com/Mheek-frosh/rexipaybackend/models"
	"github.com/gin-gonic/gin"
)

func GetBanks(c *gin.Context) {
	var banks []models.Bank
	// For now, if no banks exist, we can seed some or just return empty
	if err := config.DB.Find(&banks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch banks"})
		return
	}

	// Seed real Nigerian banks if empty
	if len(banks) == 0 {
		realBanks := []models.Bank{
			{BankCode: "044", BankName: "Access Bank"},
			{BankCode: "058", BankName: "Guaranty Trust Bank (GTB)"},
			{BankCode: "033", BankName: "United Bank for Africa (UBA)"},
			{BankCode: "057", BankName: "Zenith Bank"},
			{BankCode: "011", BankName: "First Bank of Nigeria"},
			{BankCode: "214", BankName: "First City Monument Bank (FCMB)"},
			{BankCode: "023", BankName: "Citibank Nigeria"},
			{BankCode: "050", BankName: "Ecobank Nigeria"},
			{BankCode: "070", BankName: "Fidelity Bank"},
			{BankCode: "203", BankName: "Globus Bank"},
			{BankCode: "030", BankName: "Heritage Bank"},
			{BankCode: "301", BankName: "Jaiz Bank"},
			{BankCode: "082", BankName: "Keystone Bank"},
			{BankCode: "50211", BankName: "Kuda Bank"},
			{BankCode: "999992", BankName: "OPay (Paycom)"},
			{BankCode: "999991", BankName: "PalmPay"},
			{BankCode: "101", BankName: "Providus Bank"},
			{BankCode: "221", BankName: "Stanbic IBTC Bank"},
			{BankCode: "068", BankName: "Standard Chartered Bank"},
			{BankCode: "232", BankName: "Sterling Bank"},
			{BankCode: "100", BankName: "SunTrust Bank"},
			{BankCode: "032", BankName: "Union Bank of Nigeria"},
			{BankCode: "215", BankName: "Unity Bank"},
			{BankCode: "035", BankName: "Wema Bank"},
		}
		config.DB.Create(&realBanks)
		banks = realBanks
	}

	c.JSON(http.StatusOK, banks)
}
