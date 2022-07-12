package entity

import "time"

type Otp struct {
	Id        uint `gorm:"primaryKey;autoIncrement:true"`
	CreatedAt time.Time
	DeletedAt *time.Time

	Code        string
	UserEmailId uint
}
