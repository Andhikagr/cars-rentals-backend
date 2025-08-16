package routes

import (
	"cars_rentals_backend/handlers"
	"database/sql"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
    r := mux.NewRouter()

    // Cars
    r.HandleFunc("/cars", handlers.GetCarsHandler(db)).Methods("GET")
    r.HandleFunc("/cars/{id:[0-9]+}", handlers.GetCarByIDHandler(db)).Methods("GET")

    // Bookings
    r.HandleFunc("/api/bookings", handlers.CreateBookingHandler(db)).Methods("POST")
    r.HandleFunc("/api/bookings", handlers.GetBookingsHandler(db)).Methods("GET")
    r.HandleFunc("/api/bookings/pay/{id:[0-9]+}", handlers.PayBookingHandler(db)).Methods("POST")

    // Midtrans
    r.HandleFunc("/api/bookings/snap-token", handlers.CreateSnapTransactionHandler).Methods("POST")

    return r
}

