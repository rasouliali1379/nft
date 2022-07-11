package contract

import "context"

type IPersist interface {
	Init(c context.Context) error
	Migrate(c context.Context) error
	Close(c context.Context) error
	Exists(c context.Context, entity any, conditions map[string]any) error
	Get(c context.Context, entity any, conditions map[string]any) (any, error)
	Create(c context.Context, entity any) (any, error)
	Update(c context.Context, entity any, data map[string]any) (any, error)
	Count(c context.Context, entity any, conditions map[string]any) (int, error)
}
