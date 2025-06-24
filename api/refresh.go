package api

import (
	"net/http"
	"time"

	"github.com/joac1144/bootdev-chirpy/internal/auth"
)

const RefreshPath string = "POST /api/refresh"

func (config *ApiConfig) RefreshHandler(rw http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Invalid or missing refresh token")
		return
	}

	user, err := config.Db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Failed to retrieve user from refresh token: "+err.Error())
		return
	}

	newToken, err := auth.MakeJWT(user.ID, config.Secret, time.Hour)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, "Failed to create new JWT token")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respond(rw, http.StatusOK, response{Token: newToken})
}
