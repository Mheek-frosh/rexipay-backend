package models

import "time"

type Transaction struct {
	ID               uint  `gorm:"primaryKey" json:"id"`
	SenderWalletID   *uint `json:"sender_wallet_id"`   // Pointer to allow null (e.g. deposit)
	ReceiverWalletID *uint `json:"receiver_wallet_id"` // Pointer to allow null (e.g. withdrawal/external)

	// External transfer details (stored in receiver_details JSONB in DB, but broken out here for simplicity if GORM handles it, or just mapped)
	// For simplicity, let's keep specific columns as the user schema had JSONB but we can use specific columns if GORM auto-migrates.
	// Actually, user SQL had receiver_details JSONB. Let's try to stick to specific columns for Go
	// unless we implement the JSONB Scanner/Valuer interface.
	// Let's keep the existing string fields for simplicity and compatibility with Transfer logic.
	RecipientBank string `json:"recipient_bank"`
	RecipientAcc  string `json:"recipient_account"`
	RecipientName string `json:"recipient_name"`

	Amount          float64   `gorm:"not null" json:"amount"`
	Fee             float64   `gorm:"default:0" json:"fee"`
	TransactionType string    `gorm:"not null" json:"transaction_type"` // deposit, withdrawal, transfer
	Status          string    `gorm:"default:'pending'" json:"status"`
	Reference       string    `gorm:"unique" json:"reference"`
	CreatedAt       time.Time `json:"created_at"`
}
