package models

import "gorm.io/gorm"

// Users is a struct that holds the user's name and email
type Users struct {
	gorm.Model
	ID    int64  `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}
