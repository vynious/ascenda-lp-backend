package types

import "gorm.io/gorm"

type Points struct {
	gorm.Model
	Id      string `gorm:"primaryKey"`
	UserId  string `gorm:"foreignKey:Id"`
	Balance int32
}

type GetPointsRequestBody struct{}

type GetPointsByUserRequestBody struct {
	UserId string
}

type UpdatePointsRequestBody struct {
	UserId     string
	NewBalance int32
}

type UpdatePointsResponseBody struct {
	Status bool
}
