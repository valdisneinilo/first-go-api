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
	UserDB database.DBUserInterface
}

func NewUserHandler(db database.DBUserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

// Create User godoc
// @Summary      Create User
// @Description  Criação de usuários
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      400  {object}  dto.Error
// @Failure      500  {object}  dto.Error
// @Router       /users [post]
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

	w.WriteHeader(http.StatusCreated)
}

// Get Access Token JWT  godoc
// @Summary      Get Access Token
// @Description  Access Token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body   dto.GetJwtInput  true  "user credentials"
// @Success      200	{object}  dto.GetJwtOutput
// @Failure      400  	{object}  dto.Error
// @Failure      401  	{object}  dto.Error
// @Failure      500  	{object}  dto.Error
// @Router       /users/signin [post]
func (h UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("Jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)

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

	_, token, err := jwt.Encode(map[string]interface{}{
		"sub": u.Id.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	if err != nil {
		http.Error(
			w,
			"Não foi possível gerar o token",
			http.StatusInternalServerError,
		)
		return
	}

	accessToken := dto.GetJwtOutput{Access_token: token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}
