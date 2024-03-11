package makerchecker

import "time"

type CreateMakerRequestBody struct {
	Action      MakerAction
	MakerId     string
	Description string
}

type MakerAction struct {
	Resource   string
	ActionType string
	Value      int64
	UserId     string
}

type CreateMakerResponseBody struct{}

type Transaction struct {
	Id          string
	Action      MakerAction
	MakerId     string
	Description string
	CheckerId   string
	CreatedAt   *time.Duration
	UpdatedAt   *time.Duration
}
