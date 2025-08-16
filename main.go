package main

import (
	"cars_rentals_backend/config"
	"cars_rentals_backend/handlers"
	"cars_rentals_backend/routes"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
    // Load .env dulu
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, using system env instead")
    }

    // Init DB
    db := config.InitDB()
    defer db.Close()

    // Init Midtrans Snap
    config.InitSnapClient()

    // Setup routes
    r := routes.SetupRoutes(db)

    // Setup cron untuk cleanup bookings otomatis
    c := cron.New()
    // @every 1m untuk testing, nanti bisa diubah @every 10m
    c.AddFunc("@every 1m", func() {
        log.Println("CleanupBookings jalan...")
        handlers.CleanupBookings(db)
    })
    c.Start()
    defer c.Stop()

    // Jalankan server
    log.Println("Server running at http://localhost:8080")
    err = http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal(err)
    }
}
