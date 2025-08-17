package handlers

import (
	"cars_rentals_backend/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cars_rentals_backend/models"

	"github.com/gorilla/mux"
)

func GenerateInvoiceNumber(counter int) string {
	today := time.Now().Format("20060102")
	return fmt.Sprintf("INV-%s-%04d", today, counter)
}

func CreateInvoiceHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        bookingID := params["id"]

        // already paid
        var exists bool
        err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM bookings WHERE id = ? AND status = 'paid')", bookingID).Scan(&exists)
        if err != nil || !exists {
            http.Error(w, "Booking not found or not paid", http.StatusBadRequest)
            return
        }

        // generate invoice number
        var count int
        db.QueryRow("SELECT COUNT(*) FROM invoices WHERE DATE(created_at) = CURDATE()").Scan(&count)
        invoiceNumber := GenerateInvoiceNumber(count + 1)

       
        var username, email, phone string
        var totalPrice float64
        err = db.QueryRow("SELECT username, email, phone, total_price FROM bookings WHERE id = ?", bookingID).
            Scan(&username, &email, &phone, &totalPrice)
        if err != nil {
            http.Error(w, "Failed to get booking data", http.StatusInternalServerError)
            return
        }

        // insert to invoices
        _, err = db.Exec(`
            INSERT INTO invoices (invoice_number, booking_id, username, email, phone, total_price, created_at)
            VALUES (?, ?, ?, ?, ?, ?, NOW())`,
            invoiceNumber, bookingID, username, email, phone, totalPrice,
        )
        if err != nil {
            http.Error(w, "Failed to create invoice", http.StatusInternalServerError)
            return
        }

        // delete, going to invoice
        _, _ = db.Exec("DELETE FROM bookings WHERE id = ?", bookingID)

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "message":        "Invoice created successfully",
            "invoice_number": invoiceNumber,
        })
    }
}

func GetInvoicesHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := config.DB.Query(`
        SELECT id, invoice_number, booking_id, username, email, phone, total_price, created_at 
        FROM invoices
        ORDER BY created_at DESC
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var invoices []models.Invoice
    for rows.Next() {
        var inv models.Invoice
        err := rows.Scan(&inv.ID, &inv.InvoiceNumber, &inv.BookingID, &inv.Username, &inv.Email, &inv.Phone, &inv.TotalPrice, &inv.CreatedAt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        invoices = append(invoices, inv)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(invoices)
}
