package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/valdisneinilo/first-go-api/internal/dto"
	"github.com/valdisneinilo/first-go-api/internal/entity"
	"github.com/valdisneinilo/first-go-api/internal/infra/database"
)

type UserHandler struct {
	UserDB       database.DBUserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(db database.DBUserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserDB:       db,
		Jwt:          jwt,
		JwtExpiresIn: jwtExpiresIn,
	}
}

func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user dto.CreateUserInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(
			w,
			"Não foi possível fazer o decode dos dados enviados",
			http.StatusBadRequest,
		)
		return
	}

	entity, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(
			w,
			"Não foi possível criar o usuário: \n"+err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	err = h.UserDB.CreateUser(entity)
	if err != nil {
		http.Error(
			w,
			"Internal Server Erro: \n"+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var user dto.GetJwtInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(
			w,
			"Não foi possível fazer o decode dos dados enviados",
			http.StatusBadRequest,
		)
		return
	}

	u, err := h.UserDB.FindUserByEmail(user.Email)
	if err != nil {
		http.Error(
			w,
			"Não Autorizado!",
			http.StatusUnauthorized,
		)
		return
	}

	err = u.ValidatePassword(user.Password)
	if err != nil {
		http.Error(
			w,
			"Não Autorizado!",
			http.StatusUnauthorized,
		)
		return
	}

	_, token, err := h.Jwt.Encode(map[string]interface{}{
		"sub": u.Id.String(),
		"exp": time.Now().Add(time.Second * time.Duration(h.JwtExpiresIn)).Unix(),
	})

	if err != nil {
		http.Error(
			w,
			"Não foi possível gerar o token",
			http.StatusInternalServerError,
		)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}
