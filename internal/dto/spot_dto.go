package dto

type SpotCardDTO struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Address   string  `json:"address"`
	ImageURL  string  `json:"image_url"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
