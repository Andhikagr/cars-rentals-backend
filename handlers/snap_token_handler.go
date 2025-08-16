package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var SnapClient snap.Client

// Inisialisasi client Snap
func InitSnapClient() {
    SnapClient.New("YOUR-SERVER-KEY", midtrans.Sandbox)
}

// Request body dari Flutter
type SnapRequest struct {
    OrderID    string `json:"order_id"`
    TotalPrice int64  `json:"total_price"`
}

func CreateSnapTransactionHandler(w http.ResponseWriter, r *http.Request) {
    var reqBody SnapRequest
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    snapReq := &snap.Request{
        TransactionDetails: midtrans.TransactionDetails{
            OrderID:  reqBody.OrderID,
            GrossAmt: reqBody.TotalPrice,
        },
        CreditCard: &snap.CreditCardDetails{
            Secure: true,
        },
    }

    snapToken, err := SnapClient.CreateTransactionToken(snapReq)
    if err != nil {
        log.Println("Midtrans Snap error:", err)
        http.Error(w, "Failed to create Snap token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "snap_token": snapToken,
    })
}
