package models

type Car struct {
	CarID        int    `json:"car_id"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	Image        string `json:"image"`
	Description  string `json:"description"`
	Transmission string `json:"transmission"`
	FuelType     string `json:"fuel_type"`
	Seats        int    `json:"seats"`
	PricePerDay  int    `json:"price_per_day"`
}
