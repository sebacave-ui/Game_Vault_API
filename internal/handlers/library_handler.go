package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/sebacave-ui/Game_Vault_API/internal/models"
	"github.com/sebacave-ui/Game_Vault_API/internal/services"
)

func LibraryHandler(db *sql.DB) http.HandlerFunc {
	addHandler := AddGameHandler(db)
	getHandler := GetLibraryHandler(db)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getHandler(w, r)
		case http.MethodPost:
			addHandler(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Metodo no permitido")
		}
	}
}

func LibraryByIDHandler(db *sql.DB) http.HandlerFunc {
	updateHandler := UpdateGameHandler(db)
	deleteHandler := DeleteGameHandler(db)
	statsHandler := GetLibraryStatsHandler(db)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/library/stats" {
			if r.Method == http.MethodGet {
				statsHandler(w, r)
				return
			}

			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Método no permitido")
			return
		}

		switch r.Method {
		case http.MethodPut:
			updateHandler(w, r)
		case http.MethodDelete:
			deleteHandler(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Método no permitido")
		}
	}
}

func AddGameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var game models.LibraryGame

		err := json.NewDecoder(r.Body).Decode(&game)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "JSON invalido")
			return
		}

		if game.RawgID == 0 || game.Title == "" {
			writeError(w, http.StatusBadRequest, "missing_fields", "rawg_id y title son obligatorios")
			return
		}

		err = services.AddGameToLibrary(db, game)
		if err != nil {
			if errors.Is(err, services.ErrDuplicateRawgID) {
				writeError(w, http.StatusConflict, "duplicate_game", "El juego ya existe")
				return
			}

			writeError(w, http.StatusInternalServerError, "internal_error", "Error guardando juego")
			return
		}

		writeJSON(w, http.StatusCreated, map[string]string{
			"message": "Juego agregado correctamente",
		})
	}
}

func GetLibraryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")

		games, err := services.GetLibrary(db, status)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "Error")
			return
		}

		writeJSON(w, http.StatusOK, games)
	}

}

func UpdateGameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/library/")

		if id == "" {
			writeError(w, http.StatusBadRequest, "missing_id", "ID requerido")
			return
		}

		var game models.LibraryGame
		err := json.NewDecoder(r.Body).Decode(&game)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "JSON invalido")
			return
		}

		if game.PersonalScore < 1 || game.PersonalScore > 10 {
			writeError(w, http.StatusBadRequest, "invalid_score", "personal_score debe estar entre 1 y 10")
			return
		}

		switch game.Status {
		case "pendiente", "jugando", "completado", "abandonado":
		default:
			writeError(w, http.StatusBadRequest, "invalid_status", "status invalido")
			return
		}

		err = services.UpdateLibraryGame(db, id, game)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "not_found", "Juego no encontrado")
				return
			}

			writeError(w, http.StatusInternalServerError, "internal_error", "Error actualizando juego")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Juego actualizado correctamente",
		})
	}
}

func DeleteGameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/library/")

		if id == "" {
			writeError(w, http.StatusBadRequest, "missing_id", "ID requerido")
			return
		}

		err := services.DeleteLibraryGame(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "not_found", "Juego not found")
				return
			}

			writeError(w, http.StatusInternalServerError, "internal_error", "Error eliminando juego")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Juego eliminado correctamente",
		})
	}
}

func GetLibraryStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := services.GetLibraryStats(db)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "Error obteniendo estadisticas")
			return
		}

		writeJSON(w, http.StatusOK, stats)
	}
}
