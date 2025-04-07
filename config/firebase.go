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
	// Ambil credentials dari environment variable FIREBASE_CREDENTIALS
	creds := os.Getenv("FIREBASE_CREDENTIALS")
	if creds == "" {
		log.Fatal("FIREBASE_CREDENTIALS environment variable not set")
	}

	// Inisialisasi Firebase dengan kredensial dari env
	opt := option.WithCredentialsJSON([]byte(creds))
	config := &firebase.Config{
		DatabaseURL: "https://smartvehiclesentinel-2ed68-default-rtdb.asia-southeast1.firebasedatabase.app/",
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("Error init Firebase App: %v", err)
	}

	client, err := app.Database(context.Background())
	if err != nil {
		log.Fatalf("Error connect ke Firebase DB: %v", err)
	}

	log.Println("Firebase initialized successfully")
	FirebaseClient = client
}
