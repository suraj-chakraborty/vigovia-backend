package model

type PaymentDetails struct {
	Visa           string `json:"visa"`
	Validity       string `json:"validity"`
	ProcessingDate string `json:"processingDate"`
}
