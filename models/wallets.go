package models

type Wallet struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	UserID        uint    `json:"user_id"`
	Type          string  `json:"type"` // "fiat" or "crypto"
	Currency      string  `gorm:"default:'NGN';not null" json:"currency"`
	Balance       float64 `json:"balance"`
	AccountNumber string  `gorm:"unique" json:"account_number"`
	BankName      string  `json:"bank_name"`
}
