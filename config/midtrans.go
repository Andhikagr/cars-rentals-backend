package config

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var snapClient snap.Client

func InitSnapClient() {
    serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
    snapClient.New(
        serverKey,        
        midtrans.Sandbox, 
    )
}
