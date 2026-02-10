package database

import (
	"log"
	"os"
	"status-page-monitor/internal/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := os.Getenv("DSN")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return err
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&models.Server{},
		&models.User{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migration completed")
	return nil
}
