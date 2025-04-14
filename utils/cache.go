package utils

import "sync"

var rfidRegisterCache = sync.Map{}

func SetRFIDRegister(nama string) {
	rfidRegisterCache.Store("register_mode", nama)
}

func GetRFIDRegister() (string, bool) {
	val, ok := rfidRegisterCache.Load("register_mode")
	if !ok {
		return "", false
	}
	return val.(string), true
}

func ClearRFIDRegister() {
	rfidRegisterCache.Delete("register_mode")
}
