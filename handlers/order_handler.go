package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"cars_rentals_backend/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var snapClient = &snap.Client{}

func init() {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		log.Println("MIDTRANS_SERVER_KEY kosong! Cek .env")
		return
	}
	snapClient.New(serverKey, midtrans.Sandbox)
	log.Println("MIDTRANS_SERVER_KEY=", serverKey)
}

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


	// Parse tanggal (sesuai format database)
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


	// Hitung durasi sewa
	days := int(returnDate.Sub(pickedDate).Hours() / 24)
	if days <= 0 {
		days = 1 // minimal 1 hari
	}


	// Hitung totalPrice & items
	totalPrice := 0
	items := []midtrans.ItemDetails{}

	for _, car := range booking.SelectedCars {
		price := car.PricePerDay * days
		totalPrice += price
		items = append(items, midtrans.ItemDetails{
			ID:    strconv.Itoa(car.CarID),
			Name:  car.Brand + " " + car.Model,
			Price: int64(price),
			Qty:   1,
		})
		
	}

	// Hitung driver
	if booking.SelectedDriver == "With Driver" && booking.StockDriver > 0 {
		driverPrice := 200000 * booking.StockDriver * days
		totalPrice += driverPrice
		items = append(items, midtrans.ItemDetails{
			ID:    "driver",
			Name:  "Driver",
			Price: int64(driverPrice),
			Qty:   1,
		})
	
	}



	// Buat request ke Midtrans Snap
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprintf("ORDER-%d", booking.ID),
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

	snapResp, err := snapClient.CreateTransaction(snapReq)
	if snapResp == nil || snapResp.Token == "" {

		http.Error(w, "Snap transaction failed", http.StatusInternalServerError)
		return
	}

	

	// Kirim response dengan token & totalPrice
	json.NewEncoder(w).Encode(struct {
		SnapToken  string `json:"snap_token"`
		TotalPrice int    `json:"total_price"`
	}{
		SnapToken:  snapResp.Token,
		TotalPrice: totalPrice,
	})
}
