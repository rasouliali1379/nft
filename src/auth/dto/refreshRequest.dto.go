package auth

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}