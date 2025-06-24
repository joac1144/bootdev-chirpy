package api

import (
	"encoding/json"
	"net/http"

	"github.com/joac1144/bootdev-chirpy/internal/auth"
	"github.com/joac1144/bootdev-chirpy/internal/database"
	"github.com/joac1144/bootdev-chirpy/models"
)

const CreateUserPath string = "POST /api/users"
const UpdateUserPath string = "PUT /api/users"

func (config *ApiConfig) CreateUserHandler(rw http.ResponseWriter, req *http.Request) {
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user, err := config.Db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	mappedUser := models.User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respond(rw, http.StatusCreated, mappedUser)
}

func (config *ApiConfig) UpdateUserHandler(rw http.ResponseWriter, req *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := request{}
	err := decoder.Decode(&params)
	if err != nil {
		respondError(rw, http.StatusBadRequest, "Invalid request body")
		return
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, err := auth.ValidateJWT(accessToken, config.Secret)
	if err != nil {
		respondError(rw, http.StatusUnauthorized, "Invalid token")
		return
	}

	newHashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondError(rw, http.StatusInternalServerError, "Failed to hash new password")
		return
	}

	updatedUser, err := config.Db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: newHashedPassword,
	})
	if err != nil {
		respondError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	respond(rw, http.StatusOK, models.User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})
}
