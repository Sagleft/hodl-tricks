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

type timeAPIResponse struct {
	Time string `json:"dateTime"`
}

type worldClockAPIResponse struct {
	Time string `json:"currentDateTime"`
}

type geonamesAPIResponse struct {
	Time string `json:"time"`
}
