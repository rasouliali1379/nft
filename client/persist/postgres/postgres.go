package postgres

import (
	"context"
	"errors"
	"fmt"
	"maskan/client/jtrace"
	"maskan/config"
	merror "maskan/error"
	jwt "maskan/src/jwt/entity"
	otp "maskan/src/otp/entity"
	user "maskan/src/user/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func (p *Postgres) Init(c context.Context) error {
	span, _ := jtrace.T().SpanFromContext(c, "postgres[Init]")
	defer span.Finish()

	dsn := "postgresql://" + config.C().Postgres.Username + ":" + config.C().Postgres.Password + "@" + config.C().Postgres.Host + "/" + config.C().Postgres.Schema
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	p.db = db

	return nil
}

func (p *Postgres) Migrate(c context.Context) error {
	span, _ := jtrace.T().SpanFromContext(c, "postgres[Migrate]")
	defer span.Finish()

	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&user.User{}, &jwt.Jwt{}, &otp.Otp{}); err != nil {
			return err
		}

		return nil
	})
}

func (p *Postgres) Close(c context.Context) error {
	span, _ := jtrace.T().SpanFromContext(c, "postgres[Close]")
	defer span.Finish()
	return nil
}

func (p *Postgres) Exists(c context.Context, entity any, conditions map[string]any) error {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[Exists]")
	defer span.Finish()

	tx := p.db.WithContext(ctx).Where("deleted_at is null")

	for column, value := range conditions {
		tx = tx.Where(fmt.Sprintf("%s = ?", column), value)
	}

	if err := tx.First(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merror.ErrRecordNotFound
		}

		return err
	}

	return nil
}

func (p *Postgres) Get(c context.Context, entity any, conditions map[string]any) (any, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[Get]")
	defer span.Finish()

	tx := p.db.WithContext(ctx).Where("deleted_at is null")

	for column, value := range conditions {
		tx = tx.Where(fmt.Sprintf("%s = ?", column), value)
	}

	if err := tx.Find(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (p *Postgres) Create(c context.Context, entity any) (any, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[Create]")
	defer span.Finish()

	if err := p.db.WithContext(ctx).Create(entity).Error; err != nil {
		return user.User{}, err
	}

	return entity, nil
}

func (p *Postgres) Update(c context.Context, entity any, data map[string]any) (any, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[Update]")
	defer span.Finish()

	if err := p.db.WithContext(ctx).Model(entity).Updates(data).Error; err != nil {
		return user.User{}, err
	}

	return entity, nil
}
