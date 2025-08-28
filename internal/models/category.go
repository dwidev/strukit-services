package models

type Category struct {
	BaseModel
	Name     string  `json:"name"`
	Icon     *string `json:"icon,omitempty"`
	Color    string  `json:"color"`
	IsActive *bool   `json:"isActive,omitempty"`
}

func (c Category) TableName() string {
	return "categories"
}
