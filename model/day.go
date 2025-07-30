package model

type Day struct {
	Day           int    `json:"day"`
	FormattedDate string `json:"formattedDate"`
	OriginalDate  string
	Image         []byte     `json:"image"` // JPEG/PNG bytes
	Activities    []Activity `json:"activities"`
}
