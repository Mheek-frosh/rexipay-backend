package routes

import (
	"github.com/Mheek-frosh/rexipaybackend/config"
	"github.com/Mheek-frosh/rexipaybackend/controllers"
	"github.com/Mheek-frosh/rexipaybackend/middleware"
	"github.com/Mheek-frosh/rexipaybackend/models"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// API Versioning (IMPORTANT 🔥)
	api := r.Group("/api/v1")

	// ========================
	// PUBLIC ROUTES
	// ========================
	auth := api.Group("/auth")
	{
		auth.POST("/send-otp", controllers.SendOTP)
		auth.POST("/resend-otp", controllers.ResendOTP)
		auth.POST("/verify-otp", controllers.VerifyOTP)
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// ========================
	// PROTECTED ROUTES
	// ========================
	protected := api.Group("/")
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

		// Extra features
		protected.GET("/banks", controllers.GetBanks)
		protected.POST("/beneficiaries", controllers.AddBeneficiary)
		protected.GET("/beneficiaries", controllers.GetBeneficiaries)
		protected.DELETE("/beneficiaries/:id", controllers.DeleteBeneficiary)
		protected.POST("/deposit", controllers.Deposit)
		protected.GET("/notifications", controllers.GetNotifications)
		protected.GET("/login-history", controllers.GetLoginHistory)
	}

	// ========================
	// TEST ROUTE (VERY USEFUL)
	// ========================
	api.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API v1 working ✅",
		})
	})
}
