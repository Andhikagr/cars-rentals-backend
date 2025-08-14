package handlers

import (
	"cars_rentals_backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
)

func CreateBookingHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var booking models.Booking
        if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Insert ke tabel bookings
        res, err := db.Exec(`INSERT INTO bookings 
            (username, email, phone, picked_date, return_date, selected_driver, stock_driver, street_address, district, regency, province, total_price)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
            booking.Username, booking.Email, booking.Phone, booking.PickedDate, booking.ReturnDate,
            booking.SelectedDriver, booking.StockDriver, booking.StreetAddress, booking.District,
            booking.Regency, booking.Province, booking.TotalPrice,
        )
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        bookingID, _ := res.LastInsertId()

        // Insert ke booking_details
        for _, car := range booking.SelectedCars {
            _, err := db.Exec(`INSERT INTO booking_details (booking_id, car_id) VALUES (?, ?)`, bookingID, car.CarID)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":     "success",
            "booking_id": bookingID,
        })
    }
}
