package auth

type SignUpRequest struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	NationalId     string `json:"national_id"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	LandLineNumber string `json:"land_line_number"`
	Password       string `json:"password"`
	Province       string `json:"province"`
	City           string `json:"city"`
	Address        string `json:"address"`
}
