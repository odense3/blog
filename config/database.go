package config

import (
	"blog/database/seeds"
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		cfg.PgsqlDB.DbUser,
		cfg.PgsqlDB.DbPassword,
		cfg.PgsqlDB.DbHost,
		cfg.PgsqlDB.DbPort,
		cfg.PgsqlDB.DbName,
	)

	fmt.Println("Connecting to database with:", dbConnString)
	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-1] Failed to connect to database" + cfg.PgsqlDB.DbHost)
		return nil, err
	}

	pgsql, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-2] Failed to get database connection")
		return nil, err
	}

	seeds.SeedUser(db)

	pgsql.SetMaxOpenConns(cfg.PgsqlDB.DbMaxOpen)
	pgsql.SetMaxIdleConns(cfg.PgsqlDB.DbMaxIdle)

	return &Postgres{DB: db}, nil
}
