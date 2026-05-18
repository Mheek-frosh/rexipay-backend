package models

type Bank struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	BankCode string `gorm:"unique;not null" json:"bank_code"`
	BankName string `gorm:"not null" json:"bank_name"`
}
