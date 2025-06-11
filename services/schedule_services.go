package services

import (
	"encoding/json"
	"log"
	"net/http"
	"smartVehicleSentinel/models"
)

const firebaseURLSchdeule = "https://smartvehiclesentinel-2ed68-default-rtdb.asia-southeast1.firebasedatabase.app/schedules.json"

// Firebase URL untuk jadwal
func GetActiveSchedules() ([]models.Schedule, error) {
	resp, err := http.Get(firebaseURLSchdeule)
	if err != nil {
		log.Println("❌ Gagal ambil data jadwal dari Firebase:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleMap map[string]models.Schedule
	if err := json.NewDecoder(resp.Body).Decode(&scheduleMap); err != nil {
		log.Println("❌ Gagal decode data jadwal:", err)
		return nil, err
	}

	schedules := make([]models.Schedule, 0, len(scheduleMap))
	for _, s := range scheduleMap {
		schedules = append(schedules, s)
	}

	log.Printf("📥 Jadwal dari Firebase: %+v\n", schedules)
	return schedules, nil
}

// test lokal
// func GetActiveSchedules() ([]models.Schedule, error) {
// 	file, err := os.Open("schedule.json")
// 	if err != nil {
// 		log.Println("❌ Gagal buka file jadwal:", err)
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var data map[string]models.Schedule
// 	if err := json.NewDecoder(file).Decode(&data); err != nil {
// 		log.Println("❌ Gagal decode data jadwal:", err)
// 		return nil, err
// 	}

// 	var schedules []models.Schedule
// 	for _, schedule := range data {
// 		schedules = append(schedules, schedule)
// 	}

// 	log.Printf("📥 Jadwal dari lokal: %+v\n", schedules)
// 	return schedules, nil
// }
