package entity

import "time"

type Otp struct {
	Id        string
	CreatedAt time.Time

	Code        string
	UserEmailId string
}
