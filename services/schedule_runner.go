package services

import (
	"context"
	"log"
	"smartVehicleSentinel/utils"
	"time"
)

type scheduleRuntime struct {
	startedAt  time.Time
	hasStarted bool
	hasEnded   bool
}

var runtimeMap = map[string]*scheduleRuntime{}

func StartScheduler(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Scheduler stopped.")
			return

		case <-ticker.C:
			now := utils.GetNowInWIB()

			schedules, err := GetActiveSchedules()
			if err != nil {
				log.Println("❌ Failed get schedules:", err)
				continue
			}

			currentStatus, err := GetCurrentRelayStatus()
			if err != nil {
				log.Println("❌ Gagal ambil status relay:", err)
				continue
			}

			for _, s := range schedules {
				if !s.Active {
					continue
				}

				startTime, err := time.Parse("15:04", s.StartTime)
				if err != nil {
					log.Printf("⚠️ Format waktu salah: %s", s.StartTime)
					continue
				}

				startToday := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, now.Location())

				// ⛔ Jika dinyalakan manual sebelum jadwal
				if currentStatus.LastOn.Year() == now.Year() &&
					currentStatus.LastOn.Month() == now.Month() &&
					currentStatus.LastOn.Day() == now.Day() &&
					currentStatus.LastOn.Before(startToday) {

					log.Printf("⛔ Skip jadwal %s karena kendaraan sudah dinyalakan manual sebelum waktu jadwal (LastOn: %s, Jadwal: %s)",
						s.ID, currentStatus.LastOn.Format("15:04"), startToday.Format("15:04"))
					continue
				}

				// Inisialisasi runtime jika belum ada
				if _, ok := runtimeMap[s.ID]; !ok {
					runtimeMap[s.ID] = &scheduleRuntime{}
				}

				// ✅ Eksekusi ON jika belum pernah dijalankan dan waktu sudah masuk
				if !runtimeMap[s.ID].hasStarted && now.After(startToday) && now.Sub(startToday) <= 30*time.Second {
					orderOn := []string{"contact", "key", "engine"}
					for _, t := range orderOn {
						for _, target := range s.OnTargets {
							if t == target {
								log.Printf("⚡ Menyalakan: %s", t)
								if t == "contact" && currentStatus.Contact {
									log.Printf("%s sudah ON, skip penjadwalan", t)
									continue
								}
								PublishRelayCommand(t, "on")
								UpdateRelayStatusFromCommand(t + "_on")

								if t == "contact" && contains(s.OnTargets, "engine") {
									log.Println("Menunggu 5 detik sebelum nyalakan engine...")
									time.Sleep(5 * time.Second)
								} else {
									time.Sleep(1 * time.Second)
								}
							}
						}
					}

					// Simpan waktu nyala
					runtimeMap[s.ID].startedAt = now
					runtimeMap[s.ID].hasStarted = true
					runtimeMap[s.ID].hasEnded = false
				}

				// ⏹️ Eksekusi OFF berdasarkan durasi
				if runtimeMap[s.ID].hasStarted && !runtimeMap[s.ID].hasEnded {
					endTime := runtimeMap[s.ID].startedAt.Add(time.Duration(s.DurationMinute) * time.Minute)

					if now.After(endTime) && now.Sub(endTime) <= 30*time.Second {
						orderOff := []string{"contact"}
						for _, t := range orderOff {
							for _, target := range s.OffTargets {
								if t == target {
									if t == "contact" && !currentStatus.Contact {
										log.Printf("ℹ️ %s sudah OFF, skip.", t)
										continue
									}
									log.Printf("💤 Mematikan: %s", t)
									PublishRelayCommand(t, "off")
									UpdateRelayStatusFromCommand(t + "_off")
									time.Sleep(1 * time.Second)
								}
							}
						}
						runtimeMap[s.ID].hasEnded = true
					}
				}
			}
		}
	}
}

func contains(arr []string, val string) bool {
	for _, s := range arr {
		if s == val {
			return true
		}
	}
	return false
}
