package postgres

import (
	"context"
	"errors"
	"fmt"
	"nft/client/jtrace"
	"nft/config"
	apperrors "nft/error"
	card "nft/src/card/entity"
	category "nft/src/category/entity"
	email "nft/src/email/entity"
	jwt "nft/src/jwt/entity"
	kyc "nft/src/kyc/entity"
	entity "nft/src/nft/entity"
	otp "nft/src/otp/entity"
	user "nft/src/user/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func (p *Postgres) Init(c context.Context) error {
	span, _ := jtrace.T().SpanFromContext(c, "postgres[Init]")
	defer span.Finish()

	dsn := "postgresql://" +
		config.C().Postgres.Username +
		":" + config.C().Postgres.Password +
		"@" + config.C().Postgres.Host +
		"/" + config.C().Postgres.Schema
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error happened while initializing the connection to database: %w", err)
	}

	p.db = db

	return nil
}

func (p *Postgres) Migrate(c context.Context) error {
	span, _ := jtrace.T().SpanFromContext(c, "Postgres[Migrate]")
	defer span.Finish()

	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			&category.Category{},
			&user.User{},
			&jwt.Jwt{},
			&email.Email{},
			&otp.Otp{},
			&card.Card{},
			&kyc.Kyc{},
			&entity.Nft{},
		); err != nil {
			return fmt.Errorf("error happened while migrating tables: %w", err)
		}

		return nil
	})
}

func (p *Postgres) Close(c context.Context) error {
	span, _ := jtrace.T().SpanFromContext(c, "Postgres[Close]")
	defer span.Finish()
	return nil
}

func (p *Postgres) Get(c context.Context, entity any, conditions map[string]any) (any, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "Postgres[Get]")
	defer span.Finish()

	tx := p.db.WithContext(ctx).Where("deleted_at is null")

	for column, value := range conditions {
		tx = tx.Where(fmt.Sprintf("%s = ?", column), value)
	}

	if err := tx.First(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrRecordNotFound
		}
		return nil, fmt.Errorf("error happened while searching for a record: %w", err)
	}

	return entity, nil
}

func (p *Postgres) GetAll(c context.Context, entity any, conditions map[string]any) (any, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "Postgres[GetAll]")
	defer span.Finish()

	tx := p.db.WithContext(ctx).Where("deleted_at is null")

	for column, value := range conditions {
		tx = tx.Where(fmt.Sprintf("%s = ?", column), value)
	}

	if err := tx.Find(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrRecordNotFound
		}
		return nil, fmt.Errorf("error happened while searching for a record: %w", err)
	}

	return entity, nil
}

func (p *Postgres) Create(c context.Context, entity any) (any, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "Postgres[Create]")
	defer span.Finish()

	if err := p.db.WithContext(ctx).Create(entity).Error; err != nil {
		return user.User{}, fmt.Errorf("error happened while creating a record: %w", err)
	}

	return entity, nil
}

func (p *Postgres) Update(c context.Context, entity any, data any) (any, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "Postgres[Update]")
	defer span.Finish()

	if err := p.db.WithContext(ctx).Model(entity).Updates(data).Error; err != nil {
		return user.User{}, fmt.Errorf("error happened while updating a record: %w", err)
	}

	return entity, nil
}

func (p *Postgres) Delete(c context.Context, entity any) error {
	span, ctx := jtrace.T().SpanFromContext(c, "Postgres[Delete]")
	defer span.Finish()
	if err := p.db.WithContext(ctx).Delete(entity).Error; err != nil {
		return fmt.Errorf("error happened while updating a record: %w", err)
	}
	return nil
}

func (p *Postgres) Count(c context.Context, entity any, conditions map[string]any) (int, error) {
	span, c := jtrace.T().SpanFromContext(c, "Postgres[Count]")
	defer span.Finish()

	var count int64

	tx := p.db.WithContext(c).Model(entity)

	for column, value := range conditions {
		tx = tx.Where(fmt.Sprintf("%s = ?", column), value)
	}

	if err := tx.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error happened while searching for a record: %w", err)
	}

	return int(count), nil
}

func (p *Postgres) Last(c context.Context, entity any, conditions map[string]any) (any, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "Postgres[Last]")
	defer span.Finish()

	tx := p.db.WithContext(ctx).Where("deleted_at is null")

	for column, value := range conditions {
		tx = tx.Where(fmt.Sprintf("%s = ?", column), value)
	}

	if err := tx.Last(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrRecordNotFound
		}
		return nil, fmt.Errorf("error happened while searching for a record: %w", err)
	}

	return entity, nil
}
