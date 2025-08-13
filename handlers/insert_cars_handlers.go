package handlers

import (
	"cars_rentals_backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func InsertCarHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var car models.Car
		if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		query := `
            INSERT INTO cars (brand, model, year, image, description, transmission, fuel_type, seats, price_per_day)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
            ON DUPLICATE KEY UPDATE
                year=VALUES(year),
                image=VALUES(image),
                description=VALUES(description),
								transmission=VALUES(transmission),
                fuel_type=VALUES(fuel_type),
                seats=VALUES(seats),
                price_per_day=VALUES(price_per_day);
        `

		_, err := db.Exec(query,
			car.Brand, car.Model, car.Year, car.Image, car.Description, car.Transmission, car.FuelType, car.Seats, car.PricePerDay)
		if err != nil {
			http.Error(w, "Failed to insert/update car", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Car inserted/updated successfully")
	}
}
