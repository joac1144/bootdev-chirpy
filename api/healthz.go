package api

import "net/http"

const HealthzPath string = "GET /api/healthz"

func HealthzHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write(([]byte)("OK"))
}
