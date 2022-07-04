package jwt

import "time"

type RefreshToken struct {
	Id        uint
	Token     string
	Invoked   bool
	UserId    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
