package model

import (
	"github.com/jinzhu/gorm"
)

// User is a model for user table
type User struct {
	gorm.Model
	// Email    string `json:"email"`
	// Name     string `json:"name"`
	Username string `gorm:"column:username;primary_key;not null" json:"username"`
	Email string `gorm:"column:email" json:"email"`
	PhoneNum string `gorm:"column:phone_number" json:"phone_number"`
	Password string `gorm:"column:password" json:"password"`
}
