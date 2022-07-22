package apperrors


import "errors"

var (
	ErrCardDoesntBelongToUser = errors.New("card does't belong to the user")
	ErrCardAlreadyExistsForUser = errors.New("card already exists for user")
)
