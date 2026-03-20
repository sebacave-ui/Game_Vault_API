package models

type GameResponse struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Rating    float64  `json:"rating"`
	Image     string   `json:"image"`
	Genres    []string `json:"genres"`
	Platforms []string `json:"platforms"`
}
