package models

import "time"

type User struct {
	ID               uint   `gorm:"primaryKey"`
	PhoneNumber      string `gorm:"uniqueIndex;size:255;not null"`
	RegistrationDate time.Time
}

type OTP struct {
	PhoneNumber string
	OTPCode     string
	CreatedAt   time.Time
}
