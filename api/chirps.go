package api

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/joac1144/bootdev-chirpy/internal/auth"
	"github.com/joac1144/bootdev-chirpy/internal/database"
	"github.com/joac1144/bootdev-chirpy/models"
)

const GetChirpPath string = "GET /api/chirps/{chirpId}"
const GetChirpsPath string = "GET /api/chirps"
const PostChirpsPath string = "POST /api/chirps"
const DeleteChirpPath string = "DELETE /api/chirps/{chirpId}"

func (config *ApiConfig) GetChirpHandler(rw http.ResponseWriter, req *http.Request) {
	chirpId, err := uuid.Parse(req.PathValue("chirpId"))
	if err != nil {
		respondError(rw, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := config.Db.GetChirpById(req.Context(), chirpId)
	if err != nil {
		respondError(rw, http.StatusNotFound, err.Error())
		return
	}

	mappedChirp := models.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respond(rw, http.StatusOK, mappedChirp)
}

func (config *ApiConfig) GetChirpsHandler(rw http.ResponseWriter, req *http.Request) {

	authorId := req.URL.Query().Get("author_id")
	sortBy := req.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if authorId != "" {
		authorUUID, err := uuid.Parse(authorId)
		if err != nil {
			respondError(rw, http.StatusBadRequest, "Invalid author ID")
			return
		}
		chirps, err = config.Db.GetChirpsByAuthorId(req.Context(), authorUUID)
		if err != nil {
			respondError(rw, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		chirps, err = config.Db.GetChirps(req.Context())
		if err != nil {
			respondError(rw, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if sortBy == "desc" {
		slices.SortFunc(chirps, func(a, b database.Chirp) int {
			if a.CreatedAt.Before(b.CreatedAt) {
				return 1
			} else if a.CreatedAt.After(b.CreatedAt) {
				return -1
			}
			return 0
		})
	}

	mappedChirps := make([]models.Chirp, len(chirps))
	for i, chirp := range chirps {
		mappedChirps[i] = models.Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respond(rw, http.StatusOK, mappedChirps)
}

func (config *ApiConfig) PostChirpsHandler(rw http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}
	userId, err := auth.ValidateJWT(token, config.Secret)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	type reqData struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqData{}
	err = decoder.Decode(&params)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	if len(params.Body) > 140 {
		respondError(rw, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanBody(params.Body)

	chirp, err := config.Db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userId,
	})
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	mappedChirp := models.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respond(rw, http.StatusCreated, mappedChirp)
}

func (config *ApiConfig) DeleteChirpHandler(rw http.ResponseWriter, req *http.Request) {
	chirpId, err := uuid.Parse(req.PathValue("chirpId"))
	if err != nil {
		respondError(rw, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}
	accessTokenUserId, err := auth.ValidateJWT(accessToken, config.Secret)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	chirp, err := config.Db.GetChirpById(req.Context(), chirpId)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	if chirp.UserID != accessTokenUserId {
		respondError(rw, http.StatusForbidden, "You are not allowed to delete this chirp")
		return
	}

	err = config.Db.DeleteChirpById(req.Context(), chirpId)
	if err != nil {
		respondError(rw, http.StatusNotFound, err.Error())
	}

	respond(rw, http.StatusNoContent, nil)
}

func cleanBody(input string) string {
	cleanedBody := strings.Split(input, " ")
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i := range cleanedBody {
		lowercase := strings.ToLower(cleanedBody[i])
		for _, badWord := range badWords {
			if lowercase == badWord {
				cleanedBody[i] = "****"
			}
		}
	}
	return strings.Join(cleanedBody, " ")
}
