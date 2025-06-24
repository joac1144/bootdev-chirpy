package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/joac1144/bootdev-chirpy/api"
	"github.com/joac1144/bootdev-chirpy/internal/database"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	config := &api.ApiConfig{Db: dbQueries, Platform: platform, Secret: secret, PolkaApiKey: polkaKey}

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", config.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.HandleFunc(api.HealthzPath, api.HealthzHandler)
	serveMux.HandleFunc(api.GetChirpsPath, config.GetChirpsHandler)
	serveMux.HandleFunc(api.GetChirpPath, config.GetChirpHandler)
	serveMux.HandleFunc(api.PostChirpsPath, config.PostChirpsHandler)
	serveMux.HandleFunc(api.DeleteChirpPath, config.DeleteChirpHandler)
	serveMux.HandleFunc(api.CreateUserPath, config.CreateUserHandler)
	serveMux.HandleFunc(api.UpdateUserPath, config.UpdateUserHandler)
	serveMux.HandleFunc(api.LoginPath, config.LoginHandler)
	serveMux.HandleFunc(api.RefreshPath, config.RefreshHandler)
	serveMux.HandleFunc(api.RevokePath, config.RevokeHandler)

	serveMux.HandleFunc(api.MetricsPath, config.CountHitsHandler)
	serveMux.HandleFunc(api.ResetPath, config.ResetHitsHandler)

	serveMux.HandleFunc(api.WebhooksPath, config.WebhooksHandler)

	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
