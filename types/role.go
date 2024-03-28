package types

import (
	"time"
)

type Role struct {
	Id          uint   `gorm:"primaryKey"`
	RoleName    string `gorm:"unique"`
	Users       UserList
	Permissions RolePermissionList `gorm:"foreignKey:RoleID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RolePermission struct {
	Id        uint `gorm:"primaryKey"`
	RoleID    uint `gorm:"index"`
	CanCreate bool `gorm:"default:false"`
	CanRead   bool `gorm:"default:false"`
	CanUpdate bool `gorm:"default:false"`
	CanDelete bool `gorm:"default:false"`
	Resource  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoleList []Role

type RolePermissionList []RolePermission

type CreateRoleRequestBody struct {
	RoleName    string
	Permissions *RolePermissionList
}

type GetRoleRequestBody struct {
	RoleName string
}

type DeleteRoleRequestBody struct {
	RoleName string
}

type UpdateRoleRequestBody struct {
	RoleName    string
	NewRoleName string
	Permissions *RolePermissionList
}
