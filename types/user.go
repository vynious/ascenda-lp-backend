package types

import (
	"time"
)

type User struct {
	Id        string `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Roles     []Role `gorm:"many2many:user_roles;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserList []User
