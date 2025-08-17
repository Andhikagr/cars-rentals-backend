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
    // Load .env
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, using system env instead")
    }

    // Init DB
    db := config.InitDB()
config.DB = db  // ✅ set DB global
defer db.Close()

    

    config.InitSnapClient()

    // Setup routes
    r := routes.SetupRoutes(db)

    // Setup cron 
    c := cron.New()
   
    c.AddFunc("@every 1m", func() {
        log.Println("CleanupBookings jalan...")
        handlers.CleanupBookings(db)
    })
    c.Start()
    defer c.Stop()

   
    log.Println("Server running at http://localhost:8080")
    err = http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal(err)
    }
}
