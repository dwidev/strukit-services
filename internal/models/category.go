package models

import (
	"time"

	"github.com/google/uuid"
)

var DefaultCategory = uuid.MustParse("566039d1-8937-44c3-b61d-dc069507f524")

type Category struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	Name     string  `json:"name"`
	Icon     *string `json:"icon,omitempty"`
	Color    string  `json:"color"`
	IsActive *bool   `json:"isActive,omitempty"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

func (c Category) TableName() string {
	return "categories"
}
