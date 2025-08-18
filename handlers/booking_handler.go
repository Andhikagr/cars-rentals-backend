package handlers

import (
	"cars_rentals_backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CreateBookingHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var booking models.Booking
        if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

       res, err := db.Exec(`INSERT INTO bookings
    (username, email, phone, picked_date, return_date, selected_driver, stock_driver, street_address, district, regency, province, total_price, status)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'draft')`,
    booking.Username, booking.Email, booking.Phone, booking.PickedDate, booking.ReturnDate,
    booking.SelectedDriver, booking.StockDriver, booking.StreetAddress, booking.District,
    booking.Regency, booking.Province, booking.TotalPrice,
)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        bookingID, _ := res.LastInsertId()

        for _, car := range booking.SelectedCars {
            _, err := db.Exec(`INSERT INTO booking_details (booking_id, car_id) VALUES (?, ?)`, bookingID, car.CarID)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":     "draft",
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
				   b.regency, b.province, b.total_price, b.created_at, b.status,
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
			var pickedDateStr, returnDateStr string
			var selectedDriver sql.NullString
			var stockDriver int
			var totalPrice sql.NullFloat64
			var createdAtStr sql.NullString
			var status string
			var carID sql.NullInt64
			var brand, model, image sql.NullString

			err := rows.Scan(
				&bookingID, &bUsername, &bEmail, &bPhone,
				&pickedDateStr, &returnDateStr,
				&selectedDriver, &stockDriver, &bStreet, &bDistrict,
				&bRegency, &bProvince, &totalPrice, &createdAtStr, &status,
				&carID, &brand, &model, &image,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			booking, exists := bookingsMap[bookingID]
			if !exists {
				sDriver := "Without Driver"
				if selectedDriver.Valid {
					sDriver = selectedDriver.String
				}

				tPrice := 0.0
				if totalPrice.Valid {
					tPrice = totalPrice.Float64
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
					PickedDate:     pickedDateStr,
					ReturnDate:     returnDateStr,
					SelectedDriver: sDriver,
					StockDriver:    stockDriver,
					StreetAddress:  bStreet,
					District:       bDistrict,
					Regency:        bRegency,
					Province:       bProvince,
					TotalPrice:     tPrice,
					CreatedAt:      cAt,
					Status:         status,
					SelectedCars:   []models.Car{},
				}

				bookingsMap[bookingID] = booking
			}

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

		bookings := make([]models.Booking, 0, len(bookingsMap))
		for _, b := range bookingsMap {
			bookings = append(bookings, *b)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bookings)
	}
}


//pay
func PayBookingHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        bookingID := vars["id"]

        _, err := db.Exec(`UPDATE bookings SET status='paid' WHERE id=?`, bookingID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "status": "paid",
        })
    }
}
//delete 
func CleanupBookings(db *sql.DB) {
   
    res, err := db.Exec(`
        UPDATE bookings
        SET status = 'expired', expired_at = NOW()
        WHERE status = 'draft'
          AND created_at < NOW() - INTERVAL 50 MINUTE
    `)
    if err != nil {
        log.Println("Error expiring draft bookings:", err)
    } else {
        rows, _ := res.RowsAffected()
        log.Printf("%d draft bookings expired\n", rows)
    }

    
    res, err = db.Exec(`
        DELETE FROM bookings
        WHERE status = 'expired'
          AND expired_at < NOW() - INTERVAL 1 HOUR
    `)
    if err != nil {
        log.Println("Error deleting expired bookings:", err)
    } else {
        rows, _ := res.RowsAffected()
        log.Printf("%d expired bookings deleted\n", rows)
    }

	}

	func DeleteBookingHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]

        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid id parameter", http.StatusBadRequest)
            return
        }

        res, err := db.Exec("DELETE FROM bookings WHERE id = ?", id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        rows, _ := res.RowsAffected()
        fmt.Printf("Deleted rows: %d\n", rows)

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
    }
}

// getbookingbyid
func GetBookingByID(db *sql.DB, bookingID int) (*models.Booking, error) {
    log.Println("Query booking ID:", bookingID)
    row := db.QueryRow(`
        SELECT id, username, email, phone, picked_date, return_date,
               selected_driver, stock_driver, street_address, district,
               regency, province, total_price, created_at, status
        FROM bookings
        WHERE id = ?
    `, bookingID)

    var b models.Booking
    var selectedDriver sql.NullString
    var createdAt sql.NullString
    var totalPrice sql.NullFloat64

    err := row.Scan(
        &b.ID, &b.Username, &b.Email, &b.Phone, &b.PickedDate, &b.ReturnDate,
        &selectedDriver, &b.StockDriver, &b.StreetAddress, &b.District,
        &b.Regency, &b.Province, &totalPrice, &createdAt, &b.Status,
    )
    if err != nil {
        return nil, err
    }

    if selectedDriver.Valid {
        b.SelectedDriver = selectedDriver.String
    } 

        
    if totalPrice.Valid {
    b.TotalPrice = totalPrice.Float64
}

    if createdAt.Valid && createdAt.String != "" {
        t, _ := time.Parse("2006-01-02 15:04:05", createdAt.String)
        b.CreatedAt = &t
    }

    // Ambil SelectedCars
    rows, err := db.Query(`
        SELECT c.car_id, c.brand, c.model, c.image, c.price_per_day
        FROM booking_details bd
        JOIN cars c ON bd.car_id = c.car_id
        WHERE bd.booking_id = ?
    `, bookingID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var car models.Car
        err := rows.Scan(&car.CarID, &car.Brand, &car.Model, &car.Image, &car.PricePerDay )
        if err != nil {
            return nil, err
        }
        b.SelectedCars = append(b.SelectedCars, car)
    }

    return &b, nil
}


