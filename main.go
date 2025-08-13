package main

import (
	"cars_rentals_backend/handlers"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
    dsn := "root:AG9235r*@tcp(127.0.0.1:3306)/cars_rentals"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    r := mux.NewRouter()
    r.HandleFunc("/cars", handlers.GetCarsHandler(db)).Methods("GET")
    r.HandleFunc("/cars/{id:[0-9]+}", handlers.GetCarByIDHandler(db)).Methods("GET")

    log.Println("Server running at http://localhost:8080")
    err = http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal(err)
    }
}
