package services

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/sebacave-ui/Game_Vault_API/internal/models"
)

var ErrDuplicateRawgID = errors.New("duplicate_rawg_id")

func AddGameToLibrary(db *sql.DB, game models.LibraryGame) error {
	query := `
		INSERT INTO game_library (
			rawg_id,
			title,
			genre,
			platform,
			cover_url,
			personal_note,
			personal_score,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := db.Exec(
		query,
		game.RawgID,
		game.Title,
		game.Genre,
		game.Platform,
		game.CoverURL,
		game.PersonalNote,
		game.PersonalScore,
		game.Status,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrDuplicateRawgID
		}
		return err
	}

	return nil
}

func GetLibrary(db *sql.DB, status string) ([]models.LibraryGame, error) {
	query := `
		SELECT id, rawg_id, title, genre, platform, cover_url,
		       personal_note, personal_score, status, added_at
		FROM game_library
	`

	var (
		rows *sql.Rows
		err  error
	)

	if status != "" {
		query += " WHERE status = $1 ORDER BY added_at DESC"
		rows, err = db.Query(query, status)
	} else {
		query += " ORDER BY added_at DESC"
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []models.LibraryGame

	for rows.Next() {
		var game models.LibraryGame

		err := rows.Scan(
			&game.ID,
			&game.RawgID,
			&game.Title,
			&game.Genre,
			&game.Platform,
			&game.CoverURL,
			&game.PersonalNote,
			&game.PersonalScore,
			&game.Status,
			&game.AddedAt,
		)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}

func UpdateLibraryGame(db *sql.DB, id string, game models.LibraryGame) error {
	query := `
		UPDATE game_library
		SET personal_note = $1,
		    personal_score = $2,
		    status = $3
		WHERE id = $4
	`

	result, err := db.Exec(query,
		game.PersonalNote,
		game.PersonalScore,
		game.Status,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteLibraryGame(db *sql.DB, id string) error {
	query := `DELETE FROM game_library WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func GetLibraryStats(db *sql.DB) (models.LibraryStats, error) {
	stats := models.LibraryStats{
		ByStatus: map[string]int{
			"completado": 0,
			"jugando":    0,
			"pendiente":  0,
			"abandonado": 0,
		},
	}

	querySummary := `
		SELECT 
			COUNT(*) AS total,
			COALESCE(AVG(personal_score) FILTER (WHERE personal_score BETWEEN 1 AND 10), 0)
		FROM game_library
	`

	err := db.QueryRow(querySummary).Scan(&stats.Total, &stats.AverageScore)
	if err != nil {
		return stats, err
	}

	queryStatus := `
		SELECT status, COUNT(*)
		FROM game_library
		WHERE status IS NOT NULL AND status <> ''
		GROUP BY status
	`

	rows, err := db.Query(queryStatus)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int

		err := rows.Scan(&status, &count)
		if err != nil {
			return stats, err
		}

		stats.ByStatus[status] = count
	}

	return stats, nil
}
