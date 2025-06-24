package api

import (
	"net/http"

	"github.com/joac1144/bootdev-chirpy/internal/auth"
)

const RevokePath string = "POST /api/revoke"

func (config *ApiConfig) RevokeHandler(rw http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(rw, http.StatusBadRequest, "Invalid or missing refresh token")
		return
	}

	err = config.Db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, "Failed to revoke refresh token: "+err.Error())
		return
	}

	respond(rw, http.StatusNoContent, nil)
}
