package collection

import "github.com/google/uuid"

type QueryCollection struct {
	UserId      *uuid.UUID
	CategoryIds []uuid.UUID
}
