package models

import "time"

type RevokedToken struct {
	Token     string    `gorm:"primaryKey"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
