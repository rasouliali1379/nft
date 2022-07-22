package card

type AddCardRequest struct {
	CardNumber string `json:"card_number" validate:"required"`
	IBAN       string `json:"iban" validate:"required"`
}
