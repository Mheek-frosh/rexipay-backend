package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/Mheek-frosh/rexipaybackend/config"
	"github.com/Mheek-frosh/rexipaybackend/models"
	"github.com/Mheek-frosh/rexipaybackend/services"
	"github.com/Mheek-frosh/rexipaybackend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// jwtKey is now managed in utils/jwt.go

type RegisterInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	OTPCode     string `json:"otp_code" binding:"required"`
}

type LoginInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type OTPInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyOTPInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

type KYCInput struct {
	FullName       string `json:"full_name" binding:"required"`
	Address        string `json:"address" binding:"required"`
	NIN            string `json:"nin"`
	PassportNumber string `json:"passport_number"`
	DOB            string `json:"dob" binding:"required"`
	Country        string `json:"country" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify OTP first
	var otp models.OTP
	log.Println("Input Phone:", input.PhoneNumber)
	log.Println("Input OTP:", input.OTPCode)

	if err := config.DB.Where("phone_number = ? AND code = ?", input.PhoneNumber, input.OTPCode).First(&otp).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	log.Println("Stored OTP:", otp.Code)
	log.Println("Expires At:", otp.ExpiresAt)
	log.Println("IsUsed:", otp.IsUsed)

	if otp.IsUsed || time.Now().After(otp.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP not verified or expired"})
		return
	}

	// Implicitly verify if not already verified
	if otp.VerifiedAt == nil {
		now := time.Now()
		otp.VerifiedAt = &now
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		PhoneNumber: input.PhoneNumber,
		Password:    string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Mark OTP as used
	otp.IsUsed = true
	config.DB.Save(&otp)

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "user_id": user.ID})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("phone_number = ?", input.PhoneNumber).First(&user).Error; err != nil {
		config.DB.Create(&models.LoginAttempt{
			PhoneNumber: input.PhoneNumber,
			IPAddress:   c.ClientIP(),
			Success:     false,
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		config.DB.Create(&models.LoginAttempt{
			PhoneNumber: input.PhoneNumber,
			IPAddress:   c.ClientIP(),
			Success:     false,
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Record successful login
	config.DB.Create(&models.LoginAttempt{
		PhoneNumber: input.PhoneNumber,
		IPAddress:   c.ClientIP(),
		Success:     true,
	})

	tokenString, err := utils.GenerateToken(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func SendOTP(c *gin.Context) {
	var input OTPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate 6-digit code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	otp := models.OTP{
		PhoneNumber: input.PhoneNumber,
		Code:        code,
		ExpiresAt:   time.Now().Add(time.Minute * 5), // 5 mins expiry
	}

	if err := config.DB.Create(&otp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	// Use the abstracted SMS service
	services.DefaultSMSService.SendOTP(input.PhoneNumber, code)

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func ResendOTP(c *gin.Context) {
	var input OTPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check 30s cooldown
	var lastOTP models.OTP
	if err := config.DB.Where("phone_number = ?", input.PhoneNumber).Order("created_at desc").First(&lastOTP).Error; err == nil {
		if time.Since(lastOTP.CreatedAt) < 30*time.Second {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Please wait 30 seconds before resending"})
			return
		}
	}

	SendOTP(c)
}

func VerifyOTP(c *gin.Context) {
	var input VerifyOTPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Verifying OTP - Phone:", input.PhoneNumber, "Code:", input.Code)

	var otp models.OTP
	if err := config.DB.Where("phone_number = ? AND code = ?", input.PhoneNumber, input.Code).First(&otp).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
		return
	}

	log.Println("Stored OTP Info - Code:", otp.Code, "ExpiresAt:", otp.ExpiresAt, "IsUsed:", otp.IsUsed)

	if otp.IsUsed || time.Now().After(otp.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP used or expired"})
		return
	}

	now := time.Now()
	otp.VerifiedAt = &now
	config.DB.Save(&otp)

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

func UpdateKYC(c *gin.Context) {
	var input KYCInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.FullName = input.FullName
	user.Address = input.Address
	user.NIN = input.NIN
	user.PassportNumber = input.PassportNumber
	user.DOB = input.DOB
	user.Country = input.Country

	// Mock Verification Logic
	if user.NIN != "" || user.PassportNumber != "" {
		// Simulate 3rd party API call
		user.IsVerified = true
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update KYC"})
		return
	}

	// Check if wallet exists, if not create one
	var count int64
	config.DB.Model(&models.Wallet{}).Where("user_id = ?", user.ID).Count(&count)

	if count == 0 {
		wallet := models.Wallet{
			UserID:        user.ID,
			Type:          "fiat",
			Currency:      "NGN",
			Balance:       0,
			AccountNumber: utils.GenerateAccountNumber(),
			BankName:      "Rexi Bank",
		}
		if err := config.DB.Create(&wallet).Error; err != nil {
			// Log error but don't fail the KYC request? Or maybe return warning?
			// For now, let's log and return success for KYC but failure for wallet
			log.Println("Failed to create wallet:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "KYC updated and verified", "user": user})
}

func GetLoginHistory(c *gin.Context) {
	// Since login history tracks by phone number (including failed ones), we can query by user's phone number?
	// But authenticated user has ID. We should get user phone from ID first, or just store UserID in login history (but failed logins have no ID).
	// The SQL schema has phone_number in login_attempts.
	// So we need to get current user's phone number.

	userID, _ := c.Get("user_id")
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var history []models.LoginAttempt
	if err := config.DB.Where("phone_number = ?", user.PhoneNumber).Order("created_at desc").Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch login history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
