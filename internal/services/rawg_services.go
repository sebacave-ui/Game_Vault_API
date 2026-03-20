package services

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/sebacave-ui/Game_Vault_API/internal/models"
)

func SearchGames(query string) ([]models.GameResponse, error) {

	apiKey := os.Getenv("RAWG_API_KEY")
	baseURL := os.Getenv("RAWG_BASE_URL")

	url := baseURL + "/games?key=" + apiKey + "&search=" + query

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}

	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}

	results := raw["results"].([]interface{})

	var games []models.GameResponse

	for _, r := range results {

		item := r.(map[string]interface{})

		game := models.GameResponse{}

		if item["id"] != nil {
			game.ID = int(item["id"].(float64))
		}

		if item["name"] != nil {
			game.Name = item["name"].(string)
		}

		if item["rating"] != nil {
			game.Rating = item["rating"].(float64)
		}

		if item["background_image"] != nil {
			game.Image = item["background_image"].(string)
		}

		// GENRES
		var genres []string
		if item["genres"] != nil {
			for _, g := range item["genres"].([]interface{}) {
				genre := g.(map[string]interface{})
				genres = append(genres, genre["name"].(string))
			}
		}
		game.Genres = genres

		// PLATFORMS
		var platforms []string
		if item["platforms"] != nil {
			for _, p := range item["platforms"].([]interface{}) {
				platform := p.(map[string]interface{})
				info := platform["platform"].(map[string]interface{})
				platforms = append(platforms, info["name"].(string))
			}
		}
		game.Platforms = platforms

		games = append(games, game)
	}

	return games, nil
}

func GetGameByID(id string) (models.GameResponse, error) {

	apiKey := os.Getenv("RAWG_API_KEY")
	baseURL := os.Getenv("RAWG_BASE_URL")

	url := baseURL + "/games/" + id + "?key=" + apiKey

	resp, err := http.Get(url)
	if err != nil {
		return models.GameResponse{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.GameResponse{}, err
	}

	var raw map[string]interface{}

	err = json.Unmarshal(body, &raw)
	if err != nil {
		return models.GameResponse{}, err
	}

	game := models.GameResponse{}

	game.ID = int(raw["id"].(float64))
	game.Name = raw["name"].(string)
	game.Rating = raw["rating"].(float64)
	game.Image = raw["background_image"].(string)

	var genres []string
	for _, g := range raw["genres"].([]interface{}) {
		genre := g.(map[string]interface{})
		genres = append(genres, genre["name"].(string))
	}
	game.Genres = genres

	var platforms []string
	for _, p := range raw["platforms"].([]interface{}) {
		platform := p.(map[string]interface{})
		info := platform["platform"].(map[string]interface{})
		platforms = append(platforms, info["name"].(string))
	}
	game.Platforms = platforms

	return game, nil
}
