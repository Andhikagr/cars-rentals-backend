package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// Snap client global
var snapClient snap.Client

func init() {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY not set in environment")
	}
	snapClient.New(serverKey, midtrans.Sandbox) 
}

// Request dari Flutter
type SnapRequest struct {
	TotalPrice int64  `json:"total_price"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

// Response JSON
type SnapResponse struct {
	SnapToken string `json:"snap_token"`
}

// Generate order ID unik
func generateOrderID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("ORDER-%d-%d", time.Now().UnixNano(), r.Intn(1000))
}

// Handler untuk create Snap token
func CreateSnapTransactionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqBody SnapRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Buat request Snap
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  generateOrderID(),
			GrossAmt: reqBody.TotalPrice,
		},
		CreditCard: &snap.CreditCardDetails{Secure: true},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: reqBody.Username,
			Email: reqBody.Email,
			Phone: reqBody.Phone,
		},
	}

	snapToken, err := snapClient.CreateTransactionToken(req)
	if err != nil {
		log.Println("Midtrans Snap error:", err)
		http.Error(w, "Failed to create Snap token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SnapResponse{SnapToken: snapToken})
}
