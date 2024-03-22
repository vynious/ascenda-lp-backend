package types

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id        string `gorm:"primaryKey"`
	Email     string
	FirstName string
	LastName  string
	Roles     []Role `gorm:"many2many:user_roles;"`
}
