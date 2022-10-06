package contract

import (
	"context"
	model "nft/src/talan/model"
)

type ITalanRepository interface {
	GenerateAddress(c context.Context) (*model.Address, error)
	GetBalance(c context.Context, address string) (float64, error)
	GetTransactions(c context.Context, address string) ([]model.Transaction, error)
}

type ITalanService interface {
	GenerateAddress(c context.Context) (*model.Address, error)
	GetBalance(c context.Context, m model.Talan) (float64, error)
	GetTransactions(c context.Context, m model.Talan) ([]model.Transaction, error)
}
