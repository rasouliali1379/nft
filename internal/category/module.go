package category

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewCategoryRepository),
	fx.Provide(NewCategoryService),
	fx.Provide(NewCategoryController),
)
