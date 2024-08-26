package db

import (
	"fmt"

	"github.com/malamsyah/go-skele/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres() (*gorm.DB, error) {
	conf := config.Instance()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		conf.DBHost,
		conf.DBUser,
		conf.DBPassword,
		conf.DBName,
		conf.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
