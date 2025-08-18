package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var SnapClient *snap.Client

func InitSnapClient() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env")
	}

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY kosong! Cek .env")
	}

	// SnapClient sekarang pointer dan sudah terinisialisasi benar
	SnapClient = &snap.Client{}
	SnapClient.New(serverKey, midtrans.Sandbox)
	log.Println("MIDTRANS_SERVER_KEY =", serverKey)
}
