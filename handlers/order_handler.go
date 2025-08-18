package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"cars_rentals_backend/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)




func CreateSnapTransactionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqBody struct {
		BookingID int `json:"booking_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	booking, err := GetBookingByID(config.DB, reqBody.BookingID)
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Jika token sudah ada, return token lama
	if booking.SnapToken != "" {
		json.NewEncoder(w).Encode(struct {
			SnapToken  string  `json:"snap_token"`
			TotalPrice float64 `json:"total_price"`
		}{
			SnapToken:  booking.SnapToken,
			TotalPrice: booking.TotalPrice,
		})
		return
	}

	// Parse tanggal
	pickedDate, err := time.Parse("2006-01-02 15:04:05", booking.PickedDate)
	if err != nil {
		http.Error(w, "Invalid pickedDate", http.StatusBadRequest)
		return
	}
	returnDate, err := time.Parse("2006-01-02 15:04:05", booking.ReturnDate)
	if err != nil {
		http.Error(w, "Invalid returnDate", http.StatusBadRequest)
		return
	}

	days := int(returnDate.Sub(pickedDate).Hours() / 24)
	if days <= 0 {
		days = 1
	}

	// Hitung total price & items
	totalPrice := 0.0
	items := []midtrans.ItemDetails{}

	for _, car := range booking.SelectedCars {
		price := float64(car.PricePerDay * days)
		totalPrice += price
		items = append(items, midtrans.ItemDetails{
			ID:    strconv.Itoa(car.CarID),
			Name:  car.Brand + " " + car.Model,
			Price: int64(price),
			Qty:   1,
		})
	}

	if booking.SelectedDriver == "With Driver" && booking.StockDriver > 0 {
		driverPrice := float64(200000 * booking.StockDriver * days)
		totalPrice += driverPrice
		items = append(items, midtrans.ItemDetails{
			ID:    "driver",
			Name:  "Driver",
			Price: int64(driverPrice),
			Qty:   1,
		})
	}

	// Buat OrderID unik
	orderID := fmt.Sprintf("ORDER-%d", booking.ID)

	// Siapkan request Snap
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(totalPrice),
		},
		CreditCard: &snap.CreditCardDetails{Secure: true},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: booking.Username,
			Email: booking.Email,
			Phone: booking.Phone,
			BillAddr: &midtrans.CustomerAddress{
				Address: booking.StreetAddress,
				City:    booking.Regency,
			},
		},
		Items: &items,
	}

	// Create transaction di Midtrans
	snapResp, err := config.SnapClient.CreateTransaction(snapReq)
	if snapResp == nil || snapResp.Token == "" {
		log.Println("Midtrans CreateTransaction error:", err, snapResp)
		http.Error(w, "Snap transaction failed", http.StatusInternalServerError)
		return
	}

	// Simpan Snap token ke DB
	_, err = config.DB.Exec(
		"UPDATE bookings SET snap_token=?, status='pending' WHERE id=?",
		snapResp.Token, booking.ID,
	)
	if err != nil {
		log.Println("Failed to save Snap token:", err)
	}

	// Return Snap token & total price ke Flutter
	json.NewEncoder(w).Encode(struct {
		SnapToken  string  `json:"snap_token"`
		TotalPrice float64 `json:"total_price"`
	}{
		SnapToken:  snapResp.Token,
		TotalPrice: totalPrice,
	})
}
