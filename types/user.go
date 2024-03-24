package types

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	Email     string
	FirstName string
	LastName  string
	Role      string
}
