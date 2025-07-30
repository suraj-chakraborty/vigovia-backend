package model

type ItineraryData struct {
	Name            string           `json:"name"`
	DepartureCity   string           `json:"departureCity"`
	DestinationCity string           `json:"destinationCity"`
	DepartureDate   string           `json:"departureDate"`
	ReturnDate      string           `json:"returnDate"`
	Travelers       int              `json:"travelers"`
	Flights         []Flight         `json:"flights"`
	Bookings        []Booking        `json:"bookings"`
	Payments        []Payment        `json:"payments"`
	PaymentDetails  []PaymentDetails `json:"paymentDetails"`
	Installments    []Installment    `json:"installments"`
	Activity        []Activity       `json:"activity"`
}
