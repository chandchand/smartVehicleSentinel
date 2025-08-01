package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

var FirebaseClient *db.Client

func InitFirebase() {
	ctx := context.Background()

	// prod
	// Ambil credentials dari environment variable FIREBASE_CREDENTIALS
	// credBytes := os.Getenv("FIREBASE_CREDENTIALS")
	// if credBytes == "" {
	// 	log.Fatal("FIREBASE_CREDENTIALS environment variable not set")
	// }
	// opt := option.WithCredentialsJSON([]byte(credBytes))

	// Inisialisasi Firebase dengan kredensial dari env
	credBytes, err := os.ReadFile("serviceAccountKey.json")
	if err != nil {
		log.Fatalf("Error reading credentials file: %v", err)
	}
	opt := option.WithCredentialsJSON([]byte(credBytes))
	conf := &firebase.Config{
		DatabaseURL: "https://smartvehiclesentinel-2ed68.firebaseio.com/",
	}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("Gagal inisialisasi Firebase: %v", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("Gagal koneksi ke Firebase DB: %v", err)
	}

	FirebaseClient = client
	log.Println("✅ Firebase Realtime Database berhasil dikoneksi.")
}
