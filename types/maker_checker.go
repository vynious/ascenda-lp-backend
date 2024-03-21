package types

import (
	"gorm.io/gorm"
	"time"
)

// Database Models
type Transaction struct {
	gorm.Model
	Id          string      `gorm:"type:string;primary_key;"`
	Action      MakerAction // Assuming the type of Action doesn't need to change
	MakerId     string      `gorm:"type:string;index"` // Index for query optimization
	CheckerId   string      `gorm:"type:string;index"` // Index for query optimization
	Description string
	Status      string
	Approval    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//type MakerChecker struct {
//	gorm.Model
//	MakerRoleId    string   `gorm:"type:string;"`
//	MakerRole      Role     `gorm:"foreignKey:MakerRoleId"`
//	CheckerRoleIds []string `gorm:"type:string[];"` // If using a relational DB, consider a join table
//	CheckerRoles   []Role   `gorm:"many2many:makerchecker_checker_roles;"`
//}

// Others
type CreateTransactionBody struct {
	Description  string
	ResourceType string
	ActionType   string
	RequestBody  string
	UserId       string
	//MakerId   string
}

type MakerAction struct {
	ResourceType string
	ActionType   string
	RequestBody  string
	UserId       string
}

type CreateMakerResponseBody struct {
	Txn Transaction
}

type UpdateTransactionRequestBody struct {
	TransactionId string
	Approval      bool
}

type UpdateTransactionResponseBody struct {
	Txn Transaction
}
