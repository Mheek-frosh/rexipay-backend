package models

import "time"

type Beneficiary struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	BankCode      string    `json:"bank_code"`
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_name"`
	CreatedAt     time.Time `json:"created_at"`
}
