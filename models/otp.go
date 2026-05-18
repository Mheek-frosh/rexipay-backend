package models

import "time"

type OTP struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	PhoneNumber string     `gorm:"not null" json:"phone_number"`
	Code        string     `gorm:"not null" json:"code"`
	ExpiresAt   time.Time  `json:"expires_at"`
	IsUsed      bool       `gorm:"default:false" json:"is_used"`
	VerifiedAt  *time.Time `json:"verified_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
