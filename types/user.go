package types

import (
	"time"
)

type User struct {
	Id        string `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	RoleID    *uint
	Role      *Role
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

type UserList []User

type CreateUserRequestBody struct {
	FirstName string
	LastName  string
	Email     string
	RoleName  string
}

type GetUserRequestBody struct {
	Email string
}

type DeleteUserRequestBody struct {
	Email string
}

type UpdateUserRequestBody struct {
	Email        string
	NewFirstName string
	NewLastName  string
	NewEmail     string
}
