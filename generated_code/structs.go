package main

type dataContainer struct {
	CreatedOn string `json:"createdOn"`
	UnlockOn  string `json:"unlockOn"`
	Data      []byte `json:"data"`
}

type worldTimeAPIResponse struct {
	Time string `json:"utc_datetime"`
}
