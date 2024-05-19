package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type Error struct {
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return err
	}

	return nil
}

func DateOnNextDay(day string) (time.Time, error) {

	if strings.EqualFold(day, "Today") {
		return time.Now(), nil
	} else if strings.EqualFold(day, "Tomorrow") {
		return time.Now().AddDate(0, 0, 1), nil
	}

	weekday := 0
	match := false
	for weekday < 7 {
		if strings.EqualFold(day, time.Weekday(weekday).String()) {
			match = true
			break;
		}
		weekday++
	}
	if !match {
		return time.Now(), errors.New("invalid day string")
	}

	today := time.Now()

	daysRemaining := (weekday - int(today.Weekday()))
	if daysRemaining < 0 {
		daysRemaining = 7 + daysRemaining
	}

	if daysRemaining <= 1 {
		daysRemaining = 7 - daysRemaining
	}

	return today.AddDate(0, 0, daysRemaining), nil
}