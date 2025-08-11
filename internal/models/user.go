package models

import (
	"time"
)

type UserModel struct {
	BaseModel
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

func (UserModel) TableName() string {
	return "users"
}
