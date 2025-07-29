package utils

import (
	"log"
	"time"
)

// GetNowInWIB mengembalikan waktu sekarang dalam zona Asia/Jakarta (WIB)
func GetNowInWIB() time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("Gagal memuat lokasi WIB, fallback ke UTC+7 manual.")
		fallback := time.Now().UTC().Add(7 * time.Hour)
		log.Printf("⏰ Sekarang (fallback): %s", fallback.Format(time.RFC3339))
		return fallback
	}
	now := time.Now().In(loc)
	log.Printf("⏰ Sekarang (WIB): %s", now.Format(time.RFC3339))
	return now
}
