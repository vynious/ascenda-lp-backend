package db

import (
	makerchecker "github.com/vynious/ascenda-lp-backend/types/maker-checker"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	Id          string                   `gorm:"type:string;primary_key;"`
	Action      makerchecker.MakerAction // Assuming the type of Action doesn't need to change
	MakerId     string                   `gorm:"type:string;index"` // Index for query optimization
	CheckerId   string                   `gorm:"type:string;index"` // Index for query optimization
	Description string
	Status      string
	Approval    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MakerChecker struct {
	gorm.Model
	MakerRoleId    string   `gorm:"type:string;"`
	MakerRole      Role     `gorm:"foreignKey:MakerRoleId"`
	CheckerRoleIds []string `gorm:"type:string[];"` // If using a relational DB, consider a join table
	CheckerRoles   []Role   `gorm:"many2many:makerchecker_checker_roles;"`
}

type User struct {
	gorm.Model
	Id                  string `gorm:"type:string;primary_key;"`
	Name                string
	RoleId              string `gorm:"type:string;"`
	Role                Role   `gorm:"foreignKey:RoleId"`
	Email               string
	Password            string
	MadeTransactions    []Transaction `gorm:"foreignKey:MakerId"`
	CheckedTransactions []Transaction `gorm:"foreignKey:CheckerId"`
}

type Role struct {
	gorm.Model
	Id          string `gorm:"type:string;primary_key;"`
	Name        string `gorm:"unique;not null"`
	LogAccess   Permission
	PointAccess Permission
	UserAccess  Permission
}

type Permission struct {
	Create bool
	Read   bool
	Update bool
	Delete bool
}
