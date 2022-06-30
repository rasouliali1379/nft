package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maskan/client/jtrace"
	"maskan/config"
	merror "maskan/error"
	jwt "maskan/src/jwt/entity"
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
		if err := tx.AutoMigrate(&user.User{}, &jwt.Jwt{}); err != nil {
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

func (p *Postgres) UserExists(c context.Context, columnName string, value string) error {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[Exists]")
	defer span.Finish()
	var userEntity user.User

	if err := p.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", columnName), value).First(&userEntity).Error; err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merror.ErrRecordNotFound
		}

		return err
	}

	log.Println(userEntity)

	return nil
}

func (p *Postgres) CreateUser(c context.Context, user user.User) (string, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[Create]")
	defer span.Finish()

	if err := p.db.WithContext(ctx).Create(&user).Error; err != nil {
		return "", err
	}

	return "", nil
}

func (p *Postgres) SaveToken(c context.Context, user jwt.Jwt) error {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[SaveToken]")
	defer span.Finish()

	if err := p.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateToken(c context.Context, id uint, token string) error {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[UpdateToken]")
	defer span.Finish()

	jwtEntity := jwt.Jwt{
		ID: id,
	}

	db := p.db.WithContext(ctx).First(&jwtEntity)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return merror.ErrRecordNotFound
		}

		return db.Error
	}

	jwtEntity.Token = token

	return db.Save(&jwtEntity).Error
}

func (p *Postgres) RetrieveToken(c context.Context, token string) (jwt.Jwt, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[RetrieveToken]")
	defer span.Finish()

	var record jwt.Jwt
	if err := p.db.WithContext(ctx).Where("token = ?", token).Find(&record).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return jwt.Jwt{}, merror.ErrRecordNotFound
		}

		return jwt.Jwt{}, err
	}

	return record, nil
}

func (p *Postgres) GetUser(c context.Context, columnName string, value string) (user.User, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[GetUser]")
	defer span.Finish()

	var userEntity user.User

	if err := p.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", columnName), value).Find(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.User{}, merror.ErrRecordNotFound
		}

		return user.User{}, err
	}

	return userEntity, nil
}
