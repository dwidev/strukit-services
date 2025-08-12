package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	BaseModel
	UserID       uuid.UUID      `json:"userID"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	TotalBudget  float64        `json:"totalBudget"`
	StartDate    *time.Time     `json:"startDate"`
	EndDate      *time.Time     `json:"endDate"`
	Status       *ProjectStatus `json:"status"`
	IsSoftDelete bool           `json:"isSoftDelete"`
	DeletedAt    *time.Time     `json:"deletedAt"`

	User *User `gorm:"foriegnKey:UserID;" json:"user,omitempty"`
}

type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusArchived  ProjectStatus = "archived"
	ProjectStatusDeleted   ProjectStatus = "deleted"
)

func (ps ProjectStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

func (ps *ProjectStatus) Scan(value interface{}) error {
	if value == nil {
		*ps = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		*ps = ProjectStatus(v)
	case []byte:
		*ps = ProjectStatus(v)
	default:
		return fmt.Errorf("cannot scan %T into ProjectStatus", value)
	}
	return nil
}
