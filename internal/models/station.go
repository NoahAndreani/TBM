package models

type Station struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Address        string  `json:"address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	TotalSlots     int     `json:"total_slots"`
	AvailableBikes int     `json:"available_bikes"`
	ElectricBikes  int     `json:"electric_bikes"`
	ClassicBikes   int     `json:"classic_bikes"`
	Status         string  `json:"status"` // "operational", "maintenance", "offline"
}

type StationStats struct {
	TotalRentals    int     `json:"total_rentals"`
	AverageUsage    float64 `json:"average_usage"`
	PeakHours       []int   `json:"peak_hours"`
	PopularityScore float64 `json:"popularity_score"`
}

// Calcule le taux d'occupation de la station
func (s *Station) OccupancyRate() float64 {
	if s.TotalSlots == 0 {
		return 0
	}
	return float64(s.AvailableBikes) / float64(s.TotalSlots)
}

// Vérifie si la station est disponible pour la location
func (s *Station) IsAvailableForRental() bool {
	return s.Status == "operational" && s.AvailableBikes > 0
}

// Vérifie si la station est disponible pour le retour
func (s *Station) IsAvailableForReturn() bool {
	return s.Status == "operational" && s.AvailableBikes < s.TotalSlots
}
