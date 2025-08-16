package config

import (
	"log"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var snapClient snap.Client

func InitSnapClient() {
    serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
    if serverKey == "" {
        log.Fatal("MIDTRANS_SERVER_KEY not set in .env")
    }

    snapClient.New(
        serverKey,        
        midtrans.Sandbox, 
    )
}
