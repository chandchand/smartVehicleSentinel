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
		return time.Now().UTC().Add(7 * time.Hour)
	}
	return time.Now().In(loc)
}
