package model

type Flight struct {
	DepartureDate string  `json:"departureDate"`
	FlightName    string  `json:"flightName"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	FlightPrice   float32 `json:"flightPrice"`
}
