package makerchecker

import (
	"github.com/vynious/ascenda-lp-backend/db"
)

type CreateTransactionBody struct {
	Action      MakerAction
	Description string
	//MakerId   string
}

type MakerAction struct {
	ResourceType string
	ActionType   string
	Value        int64
	UserId       string
}

type CreateMakerResponseBody struct {
	Txn db.Transaction
}

type UpdateTransactionRequestBody struct {
	TransactionId string
	Approval      bool
	//CheckerId   string
}

type UpdateTransactionResponseBody struct {
	Txn db.Transaction
}
