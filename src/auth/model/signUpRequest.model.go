package auth

type SignUpRequest struct {
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	NationalId     string `json:"national_id" validate:"required"`
	Email          string `json:"email" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required"`
	LandLineNumber string `json:"land_line_number" validate:"required"`
	Password       string `json:"password" validate:"required"`
	Province       string `json:"province" validate:"required"`
	City           string `json:"city" validate:"required"`
	Address        string `json:"address" validate:"required"`
}
