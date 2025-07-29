package models

import "time"

type Relay struct {
	Contact   bool      `json:"contact"`
	Engine    bool      `json:"engine"`
	Key       bool      `json:"key"`
	LastOn    time.Time `json:"-"`       // Waktu terakhir relay dinyalakan
	LastOnRaw any       `json:"last_on"` // Waktu terakhir relay dinyalakan dalam format string, misalnya "2025-06-11T15:48:21+07:00"
}
