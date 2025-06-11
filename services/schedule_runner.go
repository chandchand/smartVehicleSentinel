package services

import (
	"context"
	"log"
	"smartVehicleSentinel/utils"
	"time"
)

func StartScheduler(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
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
			// Ambil status relay sekarang
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
				endTime := startToday.Add(time.Duration(s.DurationMinute) * time.Minute)

				// 🔐 VALIDASI: skip schedule jika kendaraan sudah dinyalakan sebelum waktu schedule
				if currentStatus.LastOn.Year() == now.Year() &&
					currentStatus.LastOn.Month() == now.Month() &&
					currentStatus.LastOn.Day() == now.Day() &&
					currentStatus.LastOn.Before(startToday) {

					log.Printf("⛔ Skip jadwal %s karena kendaraan sudah dinyalakan manual sebelum waktu jadwal (LastOn: %s, Jadwal: %s)",
						s.ID, currentStatus.LastOn.Format("15:04"), startToday.Format("15:04"))
					continue
				}

				// log.Printf("📅 Jadwal: %s - StartToday: %s, LastOn: %s", s.ID, startToday.Format(time.RFC3339), currentStatus.LastOn.Format(time.RFC3339))

				// 🔛 Menyalakan sesuai urutan
				if now.Format("15:04") == startToday.Format("15:04") {
					orderOn := []string{"contact", "key", "engine"} // urutan pasti
					for _, t := range orderOn {
						for _, target := range s.OnTargets {
							if t == target {
								log.Printf("⚡ Menyalakan: %s", t)
								if t == "contact" && currentStatus.Contact {
									log.Printf("ℹ️ %s sudah ON, skip penjadwalan", t)
									continue
								}
								PublishRelayCommand(t, "on")            // TANPA 'go' supaya berurutan
								UpdateRelayStatusFromCommand(t + "_on") // update status
								// ⏱️ Delay 5 detik khusus sebelum "engine"
								if t == "contact" && contains(s.OnTargets, "engine") {
									log.Println("⏳ Menunggu 5 detik sebelum nyalakan engine...")
									time.Sleep(5 * time.Second)
								} else {
									time.Sleep(1 * time.Second)
								}
							}
						}
					}
				}

				// 🔻 Mematikan sesuai urutan terbalik
				if now.Format("15:04") == endTime.Format("15:04") {
					orderOff := []string{"contact"}
					for _, t := range orderOff {
						for _, target := range s.OffTargets {
							if t == target {

								if t == "contact" && !currentStatus.Contact {
									log.Printf("ℹ️ %s sudah OFF, skip.", t)
									continue
								}
								log.Printf("💤 Mematikan: %s", t)
								PublishRelayCommand(t, "off")            // TANPA 'go' supaya berurutan
								UpdateRelayStatusFromCommand(t + "_off") // update status
								time.Sleep(1 * time.Second)
							}
						}
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
