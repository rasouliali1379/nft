package user

type UserListDto struct {
	Users []UserDto `json:"users"`
}

type UserDto struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	NationalId     string `json:"national_id"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	LandLineNumber string `json:"land_line_number"`
	Province       string `json:"province"`
	City           string `json:"city"`
	Address        string `json:"address"`
	PublicKey      string `josn:"public_key"`
}
