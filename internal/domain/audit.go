package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Audit struct {
	ID        string `gorm:"type:string; primary_key; not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Audit) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.NewString()
	return nil
}
