package main

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"
)

type timeHandler struct{}

func newTimeHandler() timeHandler {
	return timeHandler{}
}

type timeParserFunc func() (*time.Time, error)

func (h *timeHandler) getCurrentTime() (*time.Time, error) {
	handlers := []timeParserFunc{
		h.parseTimeFromWorldAPI,
	}

	for i, handler := range handlers {
		timeResult, err := handler()
		if err != nil {
			log.Println("time handler #" + strconv.Itoa(i) + ": " + err.Error())
			continue
		}
		log.Println("use time result from handler #" + strconv.Itoa(i))
		return timeResult, nil
	}

	return nil, errors.New("all time api servers are offline?")
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

	timeParsed, err := time.Parse(time.RFC3339Nano, timeResult.Time)
	if err != nil {
		return nil, errors.New("failed to parse time (api): " + err.Error())
	}

	return &timeParsed, nil
}
