package models

type BookingRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	PickedDate     string `json:"pickedDate"`
	ReturnDate     string `json:"returnDate"`
	SelectedDriver string `json:"selectedDriver"`
	StockDriver    int    `json:"stockDriver"`
	StreetAddress  string `json:"streetAddress"`
	District       string `json:"district"`
	Regency        string `json:"regency"`
	Province       string `json:"province"`
	TotalPrice     int    `json:"totalPrice"`
	SelectedCars   []Car  `json:"selectedCars"`
}