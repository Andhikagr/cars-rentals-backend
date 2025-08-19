package models

import "time"

type Booking struct {
    ID             int        `json:"id"`
    Username       string     `json:"username"`
    Email          string     `json:"email"`
    Phone          string     `json:"phone"`
    PickedDate     string     `json:"pickedDate"`
    ReturnDate     string     `json:"returnDate"`
    SelectedDriver string     `json:"selectedDriver"`
    StockDriver    int        `json:"stockDriver"`
    StreetAddress  string     `json:"streetAddress"`
    District       string     `json:"district"`
    Regency        string     `json:"regency"`
    Province       string     `json:"province"`
    TotalPrice     float64        `json:"totalPrice"`
    SelectedCars   []Car      `json:"selectedCars"`
    CreatedAt      *time.Time `json:"created_at"` 
    Status         string     `json:"status"`
    SnapToken      string     `json:"snapToken"`
    PaymentType      string     `json:"paymentType"`
    TransactionId      string     `json:"transactionId"`

}
