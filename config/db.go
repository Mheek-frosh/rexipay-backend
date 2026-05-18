package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Mheek-frosh/rexipaybackend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Construct DSN from individual environment variables
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("DB_SSLMODE")

		if host == "" {
			host = "localhost"
		}
		if user == "" {
			user = "postgres"
		}
		if password == "" {
			password = "postgres"
		}
		if dbname == "" {
			dbname = "rexipaybackend"
		}
		if port == "" {
			port = "5432"
		}
		if sslmode == "" {
			sslmode = "disable"
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, password, dbname, port, sslmode)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connection established")

	// Auto-migrate the models
	err = db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.OTP{},
		&models.Transaction{},
		&models.Bank{},
		&models.Beneficiary{},
		&models.Notification{},
		&models.LoginAttempt{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = db
}
