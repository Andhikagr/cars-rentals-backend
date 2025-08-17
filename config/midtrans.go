package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var SnapClient snap.Client

func InitSnapClient() {
    // Load .env hanya sekali di sini, biar aman
    godotenv.Load()

    serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
   

    log.Println("MIDTRANS_SERVER_KEY="+serverKey)
    SnapClient.New(serverKey, midtrans.Sandbox)
}
