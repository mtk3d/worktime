package main

import (
	"errors"
	"strconv"
	"strings"
)

type Event struct {
	time   int
	action string
}

func NewEvent(eventString string) (*Event, error) {
	parts := strings.Split(eventString, "-")
	if len(parts) != 2 {
		return nil, errors.New("Parsing events error")
	}

	time, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	event := &Event{
		time:   time,
		action: parts[1],
	}

	return event, nil
}
