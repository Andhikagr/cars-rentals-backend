package handlers

import (
	"cars_rentals_backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
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
            (username, email, phone, picked_date, return_date, selected_driver, stock_driver, street_address, district, regency, province, total_price, created_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
            booking.Username, booking.Email, booking.Phone, booking.PickedDate, booking.ReturnDate,
            booking.SelectedDriver, booking.StockDriver, booking.StreetAddress, booking.District,
            booking.Regency, booking.Province, booking.TotalPrice, booking.CreatedAt,
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

//get
func GetBookingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT b.id AS booking_id,
				   b.username, b.email, b.phone, b.picked_date, b.return_date,
				   b.selected_driver, b.stock_driver, b.street_address, b.district,
				   b.regency, b.province, b.total_price, b.created_at,
				   c.car_id, c.brand, c.model, c.image
			FROM bookings b
			LEFT JOIN booking_details bd ON b.id = bd.booking_id
			LEFT JOIN cars c ON bd.car_id = c.car_id
			ORDER BY b.id DESC
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		bookingsMap := make(map[int]*models.Booking)

		for rows.Next() {
			var bookingID int
			var bUsername, bEmail, bPhone, bStreet, bDistrict, bRegency, bProvince string
			var pickedDateStr, returnDateStr string // scan sebagai string
			var selectedDriver sql.NullString
			var stockDriver int
			var totalPrice sql.NullFloat64
			var createdAtStr sql.NullString
			var carID sql.NullInt64
			var brand, model, image sql.NullString
			
			

			err := rows.Scan(
				&bookingID, &bUsername, &bEmail, &bPhone,
				&pickedDateStr, &returnDateStr, // scan sebagai string
				&selectedDriver, &stockDriver, &bStreet, &bDistrict,
				&bRegency, &bProvince, &totalPrice, &createdAtStr,
				&carID, &brand, &model, &image,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Cek apakah booking sudah ada di map
			booking, exists := bookingsMap[bookingID]
			if !exists {
				sDriver := "Without Driver"
				if selectedDriver.Valid {
					sDriver = selectedDriver.String
				}

				tPrice := 0
				if totalPrice.Valid {
					tPrice = int(totalPrice.Float64)
				}

			var cAt *time.Time
if createdAtStr.Valid && createdAtStr.String != "" {
    t, err := time.Parse("2006-01-02 15:04:05", createdAtStr.String)
    if err == nil {
        cAt = &t
    }
}


				booking = &models.Booking{
					ID:             bookingID,
					Username:       bUsername,
					Email:          bEmail,
					Phone:          bPhone,
					PickedDate:     pickedDateStr,  // langsung string
					ReturnDate:     returnDateStr,  // langsung string
					SelectedDriver: sDriver,
					StockDriver:    stockDriver,
					StreetAddress:  bStreet,
					District:       bDistrict,
					Regency:        bRegency,
					Province:       bProvince,
					TotalPrice:     tPrice,
					CreatedAt:      cAt,
					SelectedCars:   []models.Car{},
				}

				bookingsMap[bookingID] = booking
			}

			// Tambahkan mobil jika ada
			if carID.Valid && brand.Valid && model.Valid && image.Valid {
				car := models.Car{
					CarID: int(carID.Int64),
					Brand: brand.String,
					Model: model.String,
					Image: image.String,
				}
				booking.SelectedCars = append(booking.SelectedCars, car)
			}
		}

		// Convert map ke slice
		bookings := make([]models.Booking, 0, len(bookingsMap))
		for _, b := range bookingsMap {
			bookings = append(bookings, *b)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bookings)
	}
}
