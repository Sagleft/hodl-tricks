package main

import (
	"time"
)

type timeHandler struct{}

func newTimeHandler() timeHandler {
	return timeHandler{}
}

func (h *timeHandler) parseTimeFromX() (*time.Time, error) {
	// API GET
	apiURL := "http://worldtimeapi.org/api/timezone/Europe/Moscow"
	_, err := httpGET(apiURL) // responseBytes
	if err != nil {
		return nil, err
	}

	//json.Unmarshal(responseBytes)

	return nil, nil // TODO
}
