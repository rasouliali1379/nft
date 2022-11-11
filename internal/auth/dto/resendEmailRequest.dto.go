package auth

type ResendEmailRequest struct {
	Token string `json:"token"`
}
