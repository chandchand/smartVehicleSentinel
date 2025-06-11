package utils

import (
	"fmt"
	"time"
)

// GetNowInWIB mengembalikan waktu sekarang dalam zona Asia/Jakarta (WIB)
func GetNowInWIB() time.Time {
	loc := time.FixedZone("WIB", 7*60*60)
	t := time.Now().In(loc)
	fmt.Println("Waktu WIB:", t)
	return time.Now().In(loc)
}
