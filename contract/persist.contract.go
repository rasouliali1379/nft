package contract

import "context"

type IPersist interface {
	Init(c context.Context) error
	Migrate(c context.Context) error
	Close(c context.Context) error
	Get(c context.Context, entity any, conditions map[string]any) (any, error)
	GetAll(c context.Context, entity any, conditions map[string]any) (any, error)
	Create(c context.Context, entity any) (any, error)
	Update(c context.Context, entity any, data any) (any, error)
	Delete(c context.Context, entity any) error
	Count(c context.Context, entity any, conditions map[string]any) (int, error)
	Last(c context.Context, entity any, conditions map[string]any) (any, error)
}
