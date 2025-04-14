package models

import "time"

type RFIDUser struct {
	Nama      string    `json:"nama"`
	UID       string    `json:"uid"`
	Timestamp time.Time `json:"timestamp"`
}
