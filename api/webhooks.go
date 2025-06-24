package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/joac1144/bootdev-chirpy/internal/auth"
	"github.com/joac1144/bootdev-chirpy/internal/database"
)

const WebhooksPath string = "POST /api/polka/webhooks"

func (config *ApiConfig) WebhooksHandler(rw http.ResponseWriter, req *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	inputApiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Invalid API key: "+err.Error())
		return
	}

	if inputApiKey != config.PolkaApiKey {
		respondError(rw, http.StatusUnauthorized, "Invalid API key")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := request{}
	err = decoder.Decode(&params)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	if params.Event != "user.upgraded" {
		respondError(rw, http.StatusNoContent, "Unsupported event: "+params.Event)
		return
	}

	err = config.Db.UpdateChirpyRedStatus(req.Context(), database.UpdateChirpyRedStatusParams{
		ID:          params.Data.UserId,
		IsChirpyRed: true,
	})
	if err != nil {
		respondError(rw, http.StatusNotFound, "Failed to update user to Chirpy Red: "+err.Error())
		return
	}

	respond(rw, http.StatusNoContent, nil)
}
