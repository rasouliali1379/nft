package user

type UserListDto struct {
	Users []UserDto `json:"users"`
}

type UserDto struct {
	ID             string `json:"id,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	NationalId     string `json:"national_id,omitempty"`
	Email          string `json:"email,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	LandLineNumber string `json:"land_line_number,omitempty"`
	Province       string `json:"province,omitempty"`
	City           string `json:"city,omitempty"`
	Address        string `json:"address,omitempty"`
	PublicKey      string `json:"public_key,omitempty"`
}
