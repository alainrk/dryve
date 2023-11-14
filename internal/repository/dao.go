package repository

import (
	"dryve/internal/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DAO interface {
	NewFileQuery() FileQuery
	NewUserQuery() UserQuery
}

type dao struct {
	db *gorm.DB
}

func NewDAO(db *gorm.DB) DAO {
	return &dao{
		db: db,
	}
}

func NewDB(config config.DatabaseConfig) (*gorm.DB, error) {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", config.Host, config.User, config.Password, config.Database, config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}

func Automigrate(db *gorm.DB, tables []any) error {
	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return err
		}
	}
	return nil
}
