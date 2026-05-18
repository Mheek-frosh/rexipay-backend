package routes

import (
	"github.com/Mheek-frosh/rexipaybackend/config"
	"github.com/Mheek-frosh/rexipaybackend/controllers"
	"github.com/Mheek-frosh/rexipaybackend/middleware"
	"github.com/Mheek-frosh/rexipaybackend/models"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Health check / Welcome route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Welcome to RexiPay API! The server is running perfectly.",
		})
	})

	// Public routes
	auth := r.Group("/auth")
	{
		auth.POST("/send-otp", controllers.SendOTP)
		auth.POST("/resend-otp", controllers.ResendOTP)
		auth.POST("/verify-otp", controllers.VerifyOTP)
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			var user models.User
			if err := config.DB.Preload("Wallets").First(&user, userID).Error; err != nil {
				c.JSON(404, gin.H{"error": "User not found"})
				return
			}
			c.JSON(200, user)
		})

		protected.POST("/kyc", controllers.UpdateKYC)

		// Wallet routes
		protected.POST("/wallets", controllers.CreateWallet)
		protected.GET("/wallets", controllers.GetUserWallets)
		protected.GET("/wallets/transactions", controllers.GetTransactions)
		protected.POST("/wallets/lookup", controllers.LookupAccount)
		protected.POST("/wallets/transfer", controllers.TransferFunds)

		// New Features
		protected.GET("/banks", controllers.GetBanks)
		protected.POST("/beneficiaries", controllers.AddBeneficiary)
		protected.GET("/beneficiaries", controllers.GetBeneficiaries)
		protected.DELETE("/beneficiaries/:id", controllers.DeleteBeneficiary)
		protected.POST("/deposit", controllers.Deposit)
		protected.GET("/notifications", controllers.GetNotifications)
		protected.GET("/login-history", controllers.GetLoginHistory)
	}
}
