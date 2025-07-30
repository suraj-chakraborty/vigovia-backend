package model

type Booking struct {
	City       string  `json:"city"`
	CheckIn    string  `json:"checkIn"`
	CheckOut   string  `json:"checkOut"`
	Nights     int     `json:"nights"`
	HotelName  string  `json:"hotelName"`
	HotelPrice float32 `json:"HotelPrice"`
}
