package model

import "time"

type User struct {
	Id        int    `gorm:"type:int;primary_key"`
	Username  string `gorm:"type:varchar(255)"`
	Password  string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
