package main

import (
	"encoding/json"
	"time"
)

type timeHandler struct{}

func newTimeHandler() timeHandler {
	return timeHandler{}
}

func (h *timeHandler) parseTimeFromX() (*time.Time, error) {
	// API GET
	apiURL := "http://worldtimeapi.org/api/timezone/Europe/Moscow"
	responseBytes, err := httpGET(apiURL)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(responseBytes, )

	return nil, nil // TODO
}
