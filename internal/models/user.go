package models

import (
	"time"
)

type User struct {
	BaseModel
	Email           *string    `json:"email"`
	PhoneNumber     *string    `json:"phoneNumber"`
	FullName        string     `json:"fullName"`
	PasswordHash    string     `json:"passwordHash"`
	AvatarUrl       *string    `json:"avatarUrl"`
	IsActive        bool       `json:"isActive"`
	IsVerified      bool       `json:"IsVerified"`
	EmailVerifiedAt *time.Time `json:"EmailVerifiedAt"`
	LastLoginAt     *time.Time `json:"lastLoginAt"`
}

func (User) TableName() string {
	return "users"
}
