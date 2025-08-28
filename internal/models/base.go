package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	return
}

type OnlyTime time.Time

func (ot *OnlyTime) Scan(value any) error {
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
	if time.Time(ot).IsZero() {
		return nil, nil
	}

	return time.Time(ot).Format("15:04:05"), nil
}

func (ot OnlyTime) MarshalJSON() ([]byte, error) {
	if time.Time(ot).IsZero() {
		return []byte("null"), nil
	}
	return fmt.Appendf(nil, `"%s"`, ot.Format()), nil
}

func (ot *OnlyTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	str := strings.Trim(string(data), `"`)
	if str == "" {
		return nil
	}

	parsed, err := time.Parse("15:04:05", str)
	if err != nil {
		return err
	}

	*ot = OnlyTime(parsed)
	return nil
}

func (ot OnlyTime) Format() string {
	return time.Time(ot).Format("15:04:05")
}
