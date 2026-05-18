package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PhoneNumber    string         `gorm:"unique;not null" json:"phone_number"`
	Password       string         `gorm:"not null" json:"-"`
	FullName       string         `json:"full_name"`
	Address        string         `json:"address"`
	NIN            string         `gorm:"unique" json:"nin"`
	PassportNumber string         `gorm:"unique" json:"passport_number"`
	DOB            string         `json:"dob"`
	Country        string         `json:"country"`
	IsVerified     bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Wallets []Wallet `json:"wallets"` // One-to-many relationship
}
