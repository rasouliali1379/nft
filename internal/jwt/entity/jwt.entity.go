package jwt

import "time"

type Jwt struct {
	ID        uint `gorm:"primaryKey;autoIncrement:true"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	Token   string `gorm:"not null"`
	UserId  string `gorm:"not null"`
	Invoked bool   `gorm:"default:false"`
}
