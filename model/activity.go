package model

type Activity struct {
	City         string  `json:"City"`
	Activity     string  `json:"Activity"`
	Date         string  `json:"Date"`
	Time         string  `json:"Time"`
	TimeRequired string  `json:"TimeRequired"`
	Image        []byte  `json:"Image"`
	ActivityCost float32 `json:"ActivityCost"`
}
