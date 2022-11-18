package talan

import (
	"context"
	"go.uber.org/fx"
	"nft/contract"
	model "nft/internal/talan/model"
)

type TalanService struct {
	talanRepository contract.ITalanRepository
}

type TalanServiceParams struct {
	fx.In
	TalanRepository contract.ITalanRepository
}

func NewTalanService(params TalanServiceParams) contract.ITalanService {
	return &TalanService{
		talanRepository: params.TalanRepository,
	}
}

func (t TalanService) GenerateAddress(c context.Context) (*model.Address, error) {
	return t.talanRepository.GenerateAddress(c)
}

func (t TalanService) GetBalance(c context.Context, m model.Talan) (float64, error) {
	return t.talanRepository.GetBalance(c, m.Address.PublicAddress)
}

func (t TalanService) GetTransactions(c context.Context, m model.Talan) ([]model.Transaction, error) {
	return t.talanRepository.GetTransactions(c, m.Address.PublicAddress)
}
