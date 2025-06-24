package api

import "net/http"

const ResetPath string = "POST /admin/reset"

func (config *ApiConfig) ResetHitsHandler(rw http.ResponseWriter, req *http.Request) {
	if config.Platform != "dev" {
		respondError(rw, 403, "Forbidden")
		return
	}

	err := config.Db.DeleteAllUsers(req.Context())
	if err != nil {
		respondError(rw, 500, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	config.FileserverHits.Store(0)
	rw.Write(([]byte)("OK"))
}
