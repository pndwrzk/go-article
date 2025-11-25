package database

import (
	"fmt"
	"log"
	"time"

	"github.com/pndwrzk/go-article/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectPostgres() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		fmt.Println(dsn)
		log.Fatalf("gagal konek PostgreSQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(dsn)
		log.Fatalf("gagal ambil sql.DB dari GORM: %v", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	DB = db
}
