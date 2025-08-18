package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	return
}

type OnlyTime time.Time

func (ot *OnlyTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	l := "15:04:05"
	switch v := value.(type) {
	case string:
		parsed, err := time.Parse(l, v)
		if err != nil {
			return err
		}
		*ot = OnlyTime(parsed)
		return nil
	case []byte:
		parsed, err := time.Parse(l, string(v))
		if err != nil {
			return err
		}
		*ot = OnlyTime(parsed)
		return nil
	case time.Time:
		*ot = OnlyTime(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into OnlyTime", value)
	}
}

func (ot OnlyTime) Value() (driver.Value, error) {
	return time.Time(ot).Format("15:04:05"), nil
}

func (ot OnlyTime) Format() string {
	return time.Time(ot).Format("15:04:05")
}
