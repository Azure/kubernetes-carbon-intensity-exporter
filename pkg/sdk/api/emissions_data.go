package api

import (
	"time"
)

type EmissionsData struct {
	Location string    `json:"location,omitempty"`
	Time     time.Time `json:"time,omitempty"`
	Rating   float64   `json:"rating,omitempty"`
	Duration string    `json:"duration,omitempty"`
}
