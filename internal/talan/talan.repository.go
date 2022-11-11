package talan

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
	"net/url"
	"nft/config"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	dto "nft/internal/talan/dto"
	model "nft/internal/talan/model"
)

var (
	addressUrl string
)

type TalanRepository struct {
	restyClient *resty.Request
}

type TalanRepositoryParams struct {
	fx.In
}

func NewTalanRepository(params TalanRepositoryParams) contract.ITalanRepository {
	var err error
	if addressUrl, err = url.JoinPath(
		config.C().Talan.BaseUrl,
		config.C().Talan.Address); err != nil {
		return nil
	}
	return &TalanRepository{restyClient: resty.New().R().EnableTrace().SetResult(dto.Response{})}
}

func (t TalanRepository) GenerateAddress(c context.Context) (*model.Address, error) {
	span, c := jtrace.T().SpanFromContext(c, "TalanRepository[GenerateAddress]")
	defer span.Finish()

	path, err := url.JoinPath(addressUrl, config.C().Talan.Generate)
	if err != nil {
		return nil, err
	}

	response, err := t.restyClient.Get(path)
	if err != nil {
		return nil, err
	}

	talanResponse := response.Result().(*dto.Response)
	var generatedAddress dto.GeneratedAddressDto
	if talanResponse.Data != nil {
		data, err := json.Marshal(talanResponse.Data)
		if err != nil {
			return nil, apperrors.ErrUnableToParseResult
		}

		err = json.Unmarshal(data, &generatedAddress)
		if err != nil {
			return nil, apperrors.ErrUnableToParseResult
		}
	}

	return mapAddressDtoToModel(generatedAddress), nil
}

func (t TalanRepository) GetBalance(c context.Context, address string) (float64, error) {
	span, c := jtrace.T().SpanFromContext(c, "TalanRepository[GetBalance]")
	defer span.Finish()

	path, err := url.JoinPath(addressUrl, address, config.C().Talan.Balance)
	if err != nil {
		return 0, err
	}

	response, err := t.restyClient.Get(path)
	if err != nil {
		return 0, err
	}

	talanResponse := response.Result().(*dto.Response)
	if talanResponse.Error != "" {
		return 0, errors.New(talanResponse.Error)
	}

	balanceRaw, ok := talanResponse.Data.(map[string]any)["balance"]
	if !ok {
		return 0, apperrors.ErrUnableToParseResult
	}

	balance, ok := balanceRaw.(float64)
	if !ok {
		return 0, apperrors.ErrUnableToParseResult
	}
	return balance, nil
}

func (t TalanRepository) GetTransactions(c context.Context, address string) ([]model.Transaction, error) {
	span, c := jtrace.T().SpanFromContext(c, "TalanRepository[GetTransactions]")
	defer span.Finish()

	path, err := url.JoinPath(addressUrl, address, config.C().Talan.Transactions)
	if err != nil {
		return nil, err
	}

	response, err := t.restyClient.Get(path)
	if err != nil {
		return nil, err
	}

	talanResponse := response.Result().(*dto.Response)
	if talanResponse.Error != "" {
		return nil, errors.New(talanResponse.Error)
	}

	var transactions []dto.TransactionDto
	if talanResponse.Data != nil {
		data, err := json.Marshal(talanResponse.Data)
		if err != nil {
			return nil, apperrors.ErrUnableToParseResult
		}

		if err = json.Unmarshal(data, &transactions); err != nil {
			return nil, apperrors.ErrUnableToParseResult
		}
	}

	return mapTransactionDtoToModel(transactions), nil
}
