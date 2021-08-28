package main

import (
	"encoding/json"
	"errors"
	"time"
)

type timeHandler struct{}

func newTimeHandler() timeHandler {
	return timeHandler{}
}

func (h *timeHandler) parseTimeFromWorldAPI() (*time.Time, error) {
	// API GET
	apiURL := "http://worldtimeapi.org/api/timezone/Europe/Moscow"
	responseBytes, err := httpGET(apiURL)
	if err != nil {
		return nil, err
	}

	timeResult := worldTimeAPIResponse{}
	err = json.Unmarshal(responseBytes, &timeResult)
	if err != nil {
		return nil, errors.New("failed to unmarshal api response json: " + err.Error())
	}

	timeParsed, err := time.Parse("2021-08-27T23:42:50.573476+00:00", timeResult.Time)
	if err != nil {
		return nil, errors.New("failed to parse time: " + err.Error())
	}

	return &timeParsed, nil
}
