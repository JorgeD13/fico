package models

import "time"

type User struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string    `json:"name" gorm:"not null"`
	ApellidoPaterno string    `json:"apellidoPaterno" gorm:"not null"`
	ApellidoMaterno string    `json:"apellidoMaterno" gorm:"not null"`
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash    string    `json:"-" gorm:"not null"`
	CreatedAt       time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
