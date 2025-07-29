package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"smartVehicleSentinel/models"
	"strings"
	"time"
)

const firebaseURL = "https://smartvehiclesentinel-2ed68-default-rtdb.asia-southeast1.firebasedatabase.app/relay.json"

// GetRelayStatus membaca status relay dari Firebase
func GetRelayStatus() (models.Relay, error) {
	resp, err := http.Get(firebaseURL)
	if err != nil {
		log.Println("❌ Gagal ambil relay status dari Firebase:", err)
		return models.Relay{}, err
	}
	defer resp.Body.Close()

	var status models.Relay
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Println("❌ Gagal decode relay status:", err)
		return models.Relay{}, err
	}

	if lastOnStr, ok := status.LastOnRaw.(string); ok {
		t, err := time.Parse(time.RFC3339, lastOnStr)
		if err == nil {
			status.LastOn = t
		}
	}

	log.Printf("📥 Status Relay dari Firebase: %+v\n", status)
	return status, nil
}

// UpdateRelayStatus memperbarui data relay
func UpdateRelayStatus(payload map[string]interface{}) error {
	// Menyiapkan request untuk Firebase
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("❌ Error mengubah payload ke JSON:", err)
		return err
	}

	// Membuat HTTP request dengan metode PATCH
	req, err := http.NewRequest(http.MethodPatch, firebaseURL, strings.NewReader(string(jsonPayload)))
	if err != nil {
		log.Println("❌ Error membuat request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Mengirim request ke Firebase
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("❌ Error mengirim request ke Firebase:", err)
		return err
	}
	defer resp.Body.Close()

	// Cek status code dari response
	if resp.StatusCode != http.StatusOK {
		log.Println("❌ Firebase returned an error:", resp.Status)
		return fmt.Errorf("Firebase error: %v", resp.Status)
	}

	log.Println("✅ Relay status berhasil diperbarui di Firebase")
	return nil
}

func GetCurrentRelayStatus() (models.Relay, error) {
	resp, err := http.Get(firebaseURL)
	if err != nil {
		return models.Relay{}, err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return models.Relay{}, err
	}

	var relay models.Relay

	// Ambil status boolean
	if v, ok := raw["contact"].(bool); ok {
		relay.Contact = v
	}
	if v, ok := raw["engine"].(bool); ok {
		relay.Engine = v
	}
	if v, ok := raw["key"].(bool); ok {
		relay.Key = v
	}

	// Parse last_on jika string
	if v, ok := raw["last_on"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			relay.LastOn = t
		} else {
			log.Printf("⚠️ Gagal parsing last_on: %v", err)
		}
	}

	return relay, nil
}
