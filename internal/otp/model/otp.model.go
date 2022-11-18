package otp

import "time"

type Otp struct {
	Id        uint
	CreatedAt time.Time

	Code        string
	UserEmailId uint
}
