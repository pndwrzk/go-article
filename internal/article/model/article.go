package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Article struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title     string    `gorm:"size:255;not null"`
	Content   string    `gorm:"type:text;not null"`
	Photos    []Photo   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Article) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	return
}
