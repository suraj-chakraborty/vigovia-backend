package model

type Installment struct {
	Installment string `json:"installment"`
	Amount      string `json:"amount"`
	DueDate     string `json:"dueDate"`
}
