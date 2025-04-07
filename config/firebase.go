package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

var FirebaseClient *db.Client

func InitFirebase() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("serviceAccountKey.json")

	conf := &firebase.Config{
		DatabaseURL: "https://smartvehiclesentinel-2ed68-default-rtdb.asia-southeast1.firebasedatabase.app/",
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
