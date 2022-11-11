package contract

import (
	"context"
	model "nft/infra/storage/model"
	"time"
)

type IStorage interface {
	Init(c context.Context) error
	Add(c context.Context, file model.File) (string, error)
	GetUrl(c context.Context, file model.File, exp time.Duration) (string, error)
}
