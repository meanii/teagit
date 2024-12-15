package models

import "gorm.io/gorm"

// Keys is a struct that holds the private and public keys
type Keys struct {
	gorm.Model
	ID         int64  `gorm:"primaryKey"`
	PrivateKey string `gorm:"not null"`
	PublicKey  string `gorm:"not null"`
	Location   string `gorm:"not null"`
}
