package main

import (
	"encoding/json"
	"errors"
	"time"
)

// Shared test helpers used across *_test.go files.

var errTest = errors.New("test error")

func fixedDay() time.Time {
	return time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC)
}

func singleRecord() []json.RawMessage {
	return []json.RawMessage{json.RawMessage(`{"test":"data"}`)}
}
