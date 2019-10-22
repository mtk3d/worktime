package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type Stats struct {
	dailyLeft      string
	dailyPassed    string
	dailyProgress  string
	weeklyLeft     string
	weeklyPassed   string
	weeklyProgress string
}

func NewStats(events []*Event) (*Stats, error) {
	now := time.Now()
	nowSec := int(now.Unix())
	weekdayNumber := int(now.Weekday())
	lastMonday := now.AddDate(0, 0, 1-weekdayNumber)
	if weekdayNumber == 0 {
		lastMonday = now.AddDate(0, 0, -6)
	}
	lastMondayMorning := time.Date(lastMonday.Year(), lastMonday.Month(), lastMonday.Day(), 0, 0, 0, 0, lastMonday.Location())
	lastMondayMorningSec := int(lastMondayMorning.Unix())
	nowMorning := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, lastMonday.Location())
	nowMorningSec := int(nowMorning.Unix())

	daily := 0
	weekly := 0
	var lastStart int
	var lastAction string

	for _, event := range events {
		lastAction = event.action
		if event.action == "start" {
			lastStart = event.time
			continue
		}

		difference := event.time - lastStart

		if event.time > lastMondayMorningSec {
			weekly += difference
		}

		if event.time > nowMorningSec {
			daily += difference
		}
	}

	if lastAction == "start" {
		difference := nowSec - lastStart

		weekly += difference
		daily += difference
	}

	dailyProgress := math.Round(float64(daily) / (8 * 60 * 60) * 20)
	weeklyProgress := math.Round(float64(weekly) / (40 * 60 * 60) * 20)

	if dailyProgress > 20 {
		dailyProgress = 20
	}

	if weeklyProgress > 20 {
		weeklyProgress = 20
	}

	displayDailyProgress, err := generateProgress(int(dailyProgress), 20)
	if err != nil {
		return nil, err
	}

	displayWeeklyProgress, err := generateProgress(int(weeklyProgress), 20)
	if err != nil {
		return nil, err
	}

	stats := &Stats{
		dailyPassed:    timeDisplay(daily),
		dailyLeft:      timeDisplay((8 * 60 * 60) - daily),
		dailyProgress:  displayDailyProgress,
		weeklyPassed:   timeDisplay(weekly),
		weeklyLeft:     timeDisplay((40 * 60 * 60) - weekly),
		weeklyProgress: displayWeeklyProgress,
	}

	return stats, nil
}

func timeDisplay(time int) string {
	sign := ""
	if time < 0 {
		time *= -1
		sign = "+"
	}
	hours := time / 3600
	minutes := time % 3600 / 60
	seconds := time % 3600 % 60

	return fmt.Sprintf("%s%02d:%02d:%02d", sign, hours, minutes, seconds)
}

func generateProgress(progress int, max int) (string, error) {
	if progress < 0 || max < 0 {
		return "", errors.New("Progress or max can't be less than 0")
	}

	if progress > max {
		return "", errors.New("Progress can't be higher than max")
	}

	blocks := strings.Repeat("█", progress)
	dots := strings.Repeat(".", max-progress)

	return "[" + blocks + dots + "]", nil
}
