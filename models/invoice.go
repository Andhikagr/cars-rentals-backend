package models

type Invoice struct {
	ID            int     `json:"id"`
	InvoiceNumber string  `json:"invoice_number"`
	BookingID     int     `json:"booking_id"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	TotalPrice    float64 `json:"total_price"`
	CreatedAt     string  `json:"created_at"`
}
