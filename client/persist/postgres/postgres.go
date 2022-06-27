package postgres

import (
	"maskan/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func (p Postgres) Init() error {
	dsn := "postgresql://" + config.C().Postgres.Username + ":" + config.C().Postgres.Password + "@" + config.C().Postgres.Host + "/" + config.C().Postgres.Schema
	_ = "host=" + config.C().Postgres.Host + " user=" + config.C().Postgres.Username + " password=" + config.C().Postgres.Password + " dbname=" + config.C().Postgres.Schema + " port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	p.db = db

	return nil
}

func (p Postgres) Close() error {
	return nil
}
