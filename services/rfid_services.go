package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"smartVehicleSentinel/models"
	"smartVehicleSentinel/utils"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Ganti URL ini dengan path collection di Firebase kamu
const rfidFirebaseURL = "https://smartvehiclesentinel-2ed68-default-rtdb.asia-southeast1.firebasedatabase.app"

func SaveRFIDToFirebase(nama, uid string) error {
	data := models.RFIDUser{
		Nama:      nama,
		UID:       uid,
		Timestamp: time.Now(),
	}

	url := fmt.Sprintf("%s/rfidData.json", rfidFirebaseURL)

	payload, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println("❌ Error membuat request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("❌ Error mengirim request ke Firebase:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("❌ Gagal simpan RFID ke Firebase. Status:", resp.Status)
		return err
	}

	log.Println("✅ RFID berhasil disimpan ke Firebase:", uid)
	return nil
}
func SetRFIDRegister(nama string) {
	// Set register mode di Firebase jadi true
	url := fmt.Sprintf("%s/rfid_register_mode.json", rfidFirebaseURL)
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte("true")))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("❌ Gagal set register mode:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("✅ Mode daftar RFID diaktifkan untuk:", nama)
}

func HandleRFIDMessage(client mqtt.Client, msg mqtt.Message) {
	uid := string(msg.Payload())
	log.Println("📥 UID diterima:", uid)

	nama, ok := utils.GetRFIDRegister()
	if ok {
		// Simpan ke Firebase
		err := SaveRFIDToFirebase(nama, uid)
		if err != nil {
			log.Println("❌ Gagal simpan ke Firebase:", err)
			return
		}

		utils.ClearRFIDRegister()
		client.Publish("rfid/mode", 0, false, "normal")
		log.Println("✅ UID berhasil disimpan:", uid)
		log.Println("⚠️ Tidak sedang dalam mode daftar. UID diabaikan.")
		return
	}
	// Kalau bukan mode daftar, berarti cek akses
	ValidateAndTriggerAccess(uid)
}

func ValidateAndTriggerAccess(uid string) {
	url := fmt.Sprintf("%s/rfidData.json", rfidFirebaseURL)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Println("❌ Gagal cek UID:", err)
		return
	}
	defer resp.Body.Close()

	var allData map[string]struct {
		Nama      string `json:"nama"`
		Timestamp string `json:"timestamp"`
		Uid       string `json:"uid"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&allData); err != nil || len(allData) == 0 {
		log.Println("⚠️ Tidak ada data ditemukan")
		createAccessLog(uid, "UNKNOWN", "DENIED")
		return
	}

	// Cek apakah UID valid
	var nama string
	for _, v := range allData {
		if v.Uid == uid {
			nama = v.Nama
			break
		}
	}

	if nama == "" {
		log.Println("⚠️ UID tidak valid:", uid)
		createAccessLog(uid, "UNKNOWN", "DENIED")
		return
	}

	log.Println("✅ UID valid:", uid, "->", nama)

	// Ambil status contact dari Firebase
	contactStatus, err := GetRelayStatus()
	if err != nil {
		log.Println("❌ Gagal mendapatkan status contact:", err)
		return
	}

	if contactStatus.Contact {
		// Matikan semua relay
		log.Println("🔻 Relay aktif, mematikan kendaraan...")

		commands := []string{"contact_off", "key_off"}
		for _, cmd := range commands {
			utils.PublishMQTT("vehicle/relay", cmd)
			_ = UpdateRelayStatusFromCommand(cmd)
			time.Sleep(1 * time.Second)
		}
	} else {
		// Nyalakan kendaraan dengan delay sebelum engine
		log.Println("🔺 Relay mati, menyalakan kendaraan...")

		utils.PublishMQTT("vehicle/relay", "key_on")
		_ = UpdateRelayStatusFromCommand("key_on")
		time.Sleep(2 * time.Second)

		utils.PublishMQTT("vehicle/relay", "contact_on")
		_ = UpdateRelayStatusFromCommand("contact_on")
		time.Sleep(3 * time.Second) // Delay 3 detik sebelum engine

		utils.PublishMQTT("vehicle/relay", "engine_on")
		_ = UpdateRelayStatusFromCommand("engine_on")
		time.Sleep(1 * time.Second)

		createAccessLog(uid, nama, "GRANTED")
	}

}

func GetContactStatus() (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/relay/contact.json", firebaseURL)) // ambil semua relay status
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Contact bool `json:"contact"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Contact, nil
}

func createAccessLog(uid, nama, status string) {
	logData := map[string]interface{}{
		"uid":       uid,
		"nama":      nama,
		"timestamp": time.Now().Format(time.RFC3339),
		"status":    status,
	}
	payload, _ := json.Marshal(logData)

	url := fmt.Sprintf("%s/accessLog.json", rfidFirebaseURL)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("❌ Gagal simpan accessLog:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("📜 Access log ditambahkan untuk:", nama)
}
