package types

import (
	"time"
)

type Log struct {
	LogId        string
	UserId       string
	Type         string
	Action       string
	UserLocation string
	Timestamp    time.Time
	TTL          string
}

type ReturnLogData struct {
	Data []Log  `json:"data"`
	Key  string `json:"key"`
}
