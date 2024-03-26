package types

import (
	"time"

	"gorm.io/gorm"
)

// Transaction => Database Model
type Transaction struct {
	gorm.Model
	TransactionId string      `gorm:"type:string;primary_key;"`
	Action        MakerAction `gorm:"type:json"`
	MakerId       string      `gorm:"type:string;index"`
	CheckerId     string      `gorm:"type:string;default:null;index"`
	Status        string      `gorm:"type:string;default:pending"`
	Approval      bool        `gorm:"type:boolean;default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// MakerChecker => Database Model (mock)
//type MakerChecker struct {
//	gorm.Model
//	MakerRoleId    string   `gorm:"type:string;"`
//	MakerRole      Role     `gorm:"foreignKey:MakerRoleId"`
//	CheckerRoleIds []string `gorm:"type:string[];"` // If using a relational DB, consider a join table
//	CheckerRoles   []Role   `gorm:"many2many:makerchecker_checker_roles;"`
//}

// Others
type CreateTransactionBody struct {
	MakerId string      `gorm:"type:string;index"`
	Action  MakerAction `gorm:"type:json"`
}

// TODO: need to update

type MakerAction struct {
	ActionType  string      `json:"action_type"`
	RequestBody interface{} `json:"request_body"` // based off other function's request body
}

type UpdateTransactionRequestBody struct {
	MakerId       string `json:"maker_id"`
	TransactionId string `json:"transaction_id"`
	Approval      bool   `json:"approval"`
}

type TransactionResponseBody struct {
	Txn Transaction
}

type GetTransactionRequestBody struct {
	TransactionId string `json:"transaction_id"`
}
