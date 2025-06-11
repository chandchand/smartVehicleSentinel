package models

type Schedule struct {
	ID             string   `json:"id"`
	OnTargets      []string `json:"on_targets"`       // ["contact", "engine"]
	OffTargets     []string `json:"off_targets"`      // ["contact"] misalnya
	StartTime      string   `json:"start_time"`       // format "15:04"
	DurationMinute int      `json:"duration_minutes"` // durasi ON sebelum OFF otomatis
	RepeatDaily    bool     `json:"repeat_daily"`
	Active         bool     `json:"active"`
}
