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
	"time"

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

func (p *Postgres) UserExists(c context.Context, columnName string, value string) error {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[Exists]")
	defer span.Finish()
	var userEntity user.User

	err := p.db.
		WithContext(ctx).
		Where("deleted_at is null").
		Where(fmt.Sprintf("%s = ?", columnName), value).
		First(&userEntity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merror.ErrRecordNotFound
		}

		return err
	}

	return nil
}

func (p *Postgres) CreateUser(c context.Context, userEntity user.User) (user.User, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[Create]")
	defer span.Finish()

	if err := p.db.WithContext(ctx).Create(&userEntity).Error; err != nil {
		return user.User{}, err
	}

	return userEntity, nil
}

func (p *Postgres) UpdateUser(c context.Context, userEntity user.User) (user.User, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[UpdateUser]")
	defer span.Finish()

	if err := p.db.WithContext(ctx).Model(&userEntity).Updates(userEntity).Error; err != nil {
		return user.User{}, err
	}

	return userEntity, nil
}

func (p *Postgres) DeleteUser(c context.Context, userId string) error {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[DeleteUser]")
	defer span.Finish()

	var userEntity user.User
	if err := p.db.WithContext(ctx).Where("id = ?", userId).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merror.ErrRecordNotFound
		}
		return err
	}

	updateTime := time.Now()
	userEntity.DeletedAt = &updateTime

	return p.db.Save(&userEntity).Error
}

func (p *Postgres) GetUser(c context.Context, columnName string, value string) (user.User, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[GetUser]")
	defer span.Finish()

	var userEntity user.User
	err := p.db.
		WithContext(ctx).
		Where("deleted_at is null").
		Where(fmt.Sprintf("%s = ?", columnName), value).
		First(&userEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.User{}, merror.ErrRecordNotFound
		}

		return user.User{}, err
	}

	return userEntity, nil
}

func (p *Postgres) GetAllUsers(c context.Context, conditions map[string]any) ([]user.User, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "postgres[GetAllUsers]")
	defer span.Finish()

	var userList []user.User
	tx := p.db.Where("deleted_at is null").WithContext(ctx).Table("users")

	for column, value := range conditions {
		tx = tx.Where(fmt.Sprintf("%s = ?", column), value)
	}

	if err := tx.Find(&userList).Error; err != nil {
		return nil, err
	}

	return userList, nil
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
