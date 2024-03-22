package types

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Id          uint `gorm:"primaryKey"`
	RoleName    string
	Permissions []RolePermission `gorm:"foreignKey:RoleID"`
	Users       []User           `gorm:"many2many:user_roles;"`
}

type RolePermission struct {
	gorm.Model
	Id        uint `gorm:"primaryKey"`
	RoleID    uint `gorm:"index"`
	CanCreate bool `gorm:"default:false"`
	CanRead   bool `gorm:"default:false"`
	CanUpdate bool `gorm:"default:false"`
	CanDelete bool `gorm:"default:false"`
	Resource  string
}

type RoleList []Role

type RolePermissionList []RolePermission
