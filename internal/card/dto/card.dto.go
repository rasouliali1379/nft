package card

import "github.com/google/uuid"

type CardList struct {
	Cards []Card `json:"cards"`
}

type Card struct {
	ID         uuid.UUID `json:"id,omitempty"`
	CardNumber string    `json:"card_number,omitempty"`
	IBAN       string    `json:"iban,omitempty"`
	Approved   bool      `json:"approved"`
}
