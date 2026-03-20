package models

type LibraryStats struct {
	Total        int            `json:"total"`
	ByStatus     map[string]int `json:"by_status"`
	AverageScore float64        `json:"average_score"`
}
