package models

import "time"

type LoginAttempt struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PhoneNumber string    `json:"phone_number"`
	IPAddress   string    `json:"ip_address"`
	Success     bool      `json:"success"`
	CreatedAt   time.Time `json:"created_at"`
}
