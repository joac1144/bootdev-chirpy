package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/joac1144/bootdev-chirpy/internal/auth"
	"github.com/joac1144/bootdev-chirpy/internal/database"
	"github.com/joac1144/bootdev-chirpy/models"
)

const LoginPath string = "POST /api/login"

func (config *ApiConfig) LoginHandler(rw http.ResponseWriter, req *http.Request) {
	type reqData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqData{}
	err := decoder.Decode(&params)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := config.Db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, config.Secret, time.Hour)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, "Failed to create JWT token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondError(rw, http.StatusInternalServerError, "Failed to create refresh token")
		return
	}

	_, err = config.Db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondError(rw, http.StatusInternalServerError, "Failed to create refresh token in database")
		return
	}

	type response struct {
		User         models.User `json:"user"`
		AccessToken  string      `json:"token,omitempty"`
		RefreshToken string      `json:"refresh_token,omitempty"`
	}

	resp := response{
		User: models.User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		AccessToken:  token,
		RefreshToken: refreshToken,
	}
	respond(rw, http.StatusOK, resp)
}
