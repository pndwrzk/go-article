package database

import (
	"log"

	"github.com/pndwrzk/go-article/internal/article/model"
)

func Migrate() {
	err := DB.AutoMigrate(
		&model.Article{},
		&model.Photo{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
