package models

import "gorm.io/gorm"

type Configs struct {
	gorm.Model
	ID     int64  `gorm:"primaryKey"`
	Host   string `gorm:"not null"`
	UserId int64  `gorm:"not null"`
	KeyId  int64  `gorm:"not null"`
}
