package model

type Payment struct {
	Label string   `json:"label"`
	Value []string `json:"value"`
}
