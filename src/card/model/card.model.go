package card

import "github.com/google/uuid"

type Card struct {
	ID         uuid.UUID
	CardNumber string
	IBAN       string
	UserId     uuid.UUID
	ApprovedBy *uuid.UUID
}
