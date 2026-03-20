package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/sebacave-ui/Game_Vault_API/internal/database"
	"github.com/sebacave-ui/Game_Vault_API/internal/handlers"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	db := database.Connect()
	defer db.Close()

	http.HandleFunc("/api/search", handlers.SearchGamesHandler)
	http.HandleFunc("/api/games/", handlers.GetGameByIDHandler)
	http.HandleFunc("/api/library", handlers.LibraryHandler(db))
	http.HandleFunc("/api/library/", handlers.LibraryByIDHandler(db))

	port := os.Getenv("PORT")

	log.Println("Servidor corriendo en el puerto", port)

	http.ListenAndServe(":"+port, nil)
}
