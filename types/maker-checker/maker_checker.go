package makerchecker

import "time"

type CreateTransactionBody struct {
	Action MakerAction
	//MakerId     string
	Description string
}

type MakerAction struct {
	ResourceType string
	ActionType   string
	Value        int64
	UserId       string
}

type CreateMakerResponseBody struct {
	Txn Transaction
}

type Transaction struct {
	Id          string
	Action      MakerAction
	MakerId     string
	Description string
	CheckerId   string
	Status      string
	Approval    bool
	CreatedAt   *time.Duration
	UpdatedAt   *time.Duration
}

type UpdateTransactionRequestBody struct {
	TransactionId string
	//CheckerId     string
	Approval bool
}

type UpdateTransactionResponseBody struct {
	Txn Transaction
}

type Email struct {
}
