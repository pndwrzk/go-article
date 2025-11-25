package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ArticleID uuid.UUID `gorm:"type:uuid;index;not null"`
	Path      string    `gorm:"size:500;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
