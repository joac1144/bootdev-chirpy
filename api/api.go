package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/joac1144/bootdev-chirpy/internal/database"
)

type ApiConfig struct {
	Db             *database.Queries
	FileserverHits atomic.Int32
	Platform       string
	Secret         string
	PolkaApiKey    string
}

func (config *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		config.FileserverHits.Add(1)
		next.ServeHTTP(rw, req)
	})
}

func respond(rw http.ResponseWriter, statusCode int, payload any) {
	rw.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		rw.WriteHeader(500)
		return
	}
	rw.WriteHeader(statusCode)
	rw.Write(data)
}

func respondError(rw http.ResponseWriter, statusCode int, errMsg string) {
	type resError struct {
		Error string `json:"error"`
	}

	respond(rw, statusCode, resError{Error: errMsg})
}
