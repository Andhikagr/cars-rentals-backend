package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func MidtransNotificationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Parse payload JSON dari Midtrans
		var notif struct {
			TransactionStatus string `json:"transaction_status"`
			OrderID           string `json:"order_id"`
			GrossAmount       string `json:"gross_amount"`
			PaymentType			string `json:"payment_type"`  
			TransactionID		string `json:"transaction_id"`  
		}
		if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
			log.Println("Failed to decode Midtrans notification:", err)
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		log.Println("Midtrans notification received:", notif)

		// Ambil bookingID dari orderID (misal ORDER-18 → 18)
		orderParts := strings.Split(notif.OrderID, "-")
		if len(orderParts) != 2 {
			http.Error(w, "Invalid order_id format", http.StatusBadRequest)
			return
		}
		bookingID := orderParts[1]

		// Cari booking di database
		var totalPrice float64
		err := db.QueryRow(`SELECT total_price FROM bookings WHERE id=?`, bookingID).Scan(&totalPrice)
		if err != nil {
			log.Println("Booking not found:", err)
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}

		// Validasi gross_amount
		grossAmt, _ := strconv.ParseFloat(notif.GrossAmount, 64)
		if grossAmt != totalPrice {
			log.Println("Gross amount mismatch:", grossAmt, "vs DB:", totalPrice)
			http.Error(w, "Gross amount mismatch", http.StatusBadRequest)
			return
		}

		// Tentukan status baru
		var newStatus string
		var updatePaidAt bool
		switch notif.TransactionStatus {
		case "settlement":
			newStatus = "paid"
			updatePaidAt = true
		case "pending":
			newStatus = "pending"
		case "deny", "cancel", "failure":
			newStatus = "failed"
		default:
			newStatus = "unknown"
		}

		// Update status booking di database
		if updatePaidAt {
			_, err = db.Exec(`UPDATE bookings SET status=?, paid_at=?, payment_type=?, transaction_id=? WHERE id=?`, 
				newStatus, 
				time.Now(),
				notif.PaymentType,
				notif.TransactionID, 
				bookingID,
			)
		} else {
			_, err = db.Exec(`UPDATE bookings SET status=? WHERE id=?`, 
			newStatus,  
			bookingID)
		}
		if err != nil {
			log.Println("Failed to update booking status:", err)
			http.Error(w, "Failed to update status", http.StatusInternalServerError)
			return
		}

		log.Println("Booking", bookingID, "status updated to", newStatus)

		// Kirim response 200 OK ke Midtrans
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
