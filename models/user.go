package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	UserName  string `json:"username" gorm:"unique"`
	Password  string `json:"password"`
	Email     string `json:"email" gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
