package models

import (
	"strukit-services/internal/models"
	"time"
)

type User struct {
	models.Base
	Email        *string    `json:"email"`
	PhoneNumber  *string    `json:"phoneNumber"`
	FullName     string     `json:"fullName"`
	PasswordHas  string     `json:"passwordHas"`
	AvatarUrl    *string    `json:"avatarUrl"`
	IsActive     bool       `json:"isActive"`
	IsVerified   bool       `json:"IsVerified"`
	IsVerifiedAt *time.Time `json:"isVerifiedAt"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`
}

func (User) TableName() string {
	return "users"
}
