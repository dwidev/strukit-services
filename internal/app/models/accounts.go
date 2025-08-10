package models

import "strukit-services/internal/models"

type Account struct {
	models.Base
	username    string  `gorm:"type:varchar(255);column:full_name" json:"fullName"`
	Email       *string `gorm:"type:varchar(255);unique" json:"email"`
	PhoneNumber *string `gorm:"type:varchar(50);column:phone_number" json:"phoneNUmber"`
	AvatarUrl   *string `gorm:"type:text;column:avatar_url" json:"avatarUrl"`
	IsActive    bool    `gorm:"type:boolean;default:true;column:is_active;" json:"isActive"`
}

func (Account) TableName() string {
	return "accounts"
}
