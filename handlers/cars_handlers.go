package handlers

import (
	"cars_rentals_backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)



func GetCarsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        search := r.URL.Query().Get("search") 

        var rows *sql.Rows
        var err error

        if search != "" {
            query := `
                SELECT car_id, brand, model, year, image, description, transmission, fuel_type, seats, price_per_day
                FROM cars
                WHERE brand LIKE ? OR model LIKE ? OR transmission LIKE ? OR seats LIKE ?`
            rows, err = db.Query(query, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
        } else {
            rows, err = db.Query("SELECT car_id, brand, model, year, image, description, transmission, fuel_type, seats, price_per_day FROM cars")
        }

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var cars []models.Car
        for rows.Next() {
            var car models.Car
            if err := rows.Scan(&car.CarID, &car.Brand, &car.Model, &car.Year, &car.Image, &car.Description, &car.Transmission, &car.FuelType, &car.Seats, &car.PricePerDay); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            cars = append(cars, car)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cars)
    }
}


//get by id
func GetCarByIDHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        
        idStr := r.URL.Path[len("/cars/"):]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid car ID", http.StatusBadRequest)
            return
        }

        var car models.Car
        query := `SELECT car_id, brand, model, year, image, description, transmission, fuel_type, seats, price_per_day FROM cars WHERE car_id = ?`
        err = db.QueryRow(query, id).Scan(&car.CarID, &car.Brand, &car.Model, &car.Year, &car.Image, &car.Description, &car.Transmission, &car.FuelType, &car.Seats, &car.PricePerDay)
        if err == sql.ErrNoRows {
            http.Error(w, "Car not found", http.StatusNotFound)
            return
        } else if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(car)
    }
}
