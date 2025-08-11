package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *BaseModel) BeforeCreate() {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
}
