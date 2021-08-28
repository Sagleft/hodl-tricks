package main

import "time"

type dataContainer struct {
	CreatedOn string `json:"createdOn"`
	UnlockOn  string `json:"unlockOn"`
	Data      []byte `json:"data"`
}

type encryptResult struct {
	TimeTo time.Time
}

type worldTimeAPIResponse struct {
	Time string `json:"utc_datetime"`
}
