package handlers

import (
	"net/http"
	"strings"

	"github.com/sebacave-ui/Game_Vault_API/internal/services"
)

func SearchGamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Método no permitido")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "missing_query", "Falta el parámetro q")
		return
	}

	result, err := services.SearchGames(query)
	if err != nil {
		writeError(w, http.StatusBadGateway, "rawg_error", "Error contacting RAWG API")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func GetGameByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Método no permitido")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/games/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing_game_id", "Game ID required")
		return
	}

	result, err := services.GetGameByID(id)
	if err != nil {
		writeError(w, http.StatusBadGateway, "rawg_error", "Error contacting RAWG API")
		return
	}

	writeJSON(w, http.StatusOK, result)
}
