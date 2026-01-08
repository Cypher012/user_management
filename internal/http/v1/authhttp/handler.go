package authhttp

import (
	"errors"
	"net/http"

	"github.com/Cypher012/userauth/internal/auth"
	"github.com/Cypher012/userauth/internal/http/httputil"
)

type AuthHandler struct {
	service *auth.AuthService
	jwt     *auth.JWTAuth
}

func NewAuthHandler(service *auth.AuthService, jwt *auth.JWTAuth) *AuthHandler {
	return &AuthHandler{
		service: service,
		jwt:     jwt,
	}
}

func (h *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req UserAuthRequest

	if !httputil.DecodeJSONBody[UserAuthRequest](w, r, &req) {
		return
	}

	if req.Email == "" || req.Password == "" {
		httputil.ErrorResponse(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := h.service.RegisterUser(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserAlreadyExists):
			httputil.ErrorResponse(w, http.StatusConflict, err.Error())
		default:
			httputil.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	atkToken, rtkToken, err := h.jwt.GenerateToken(user.ID)
	if err != nil {
		httputil.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	payload := UserResponse{
		Message: "User sign up succesful",
		Auth: Auth{
			Atk: atkToken,
			Rtk: rtkToken,
		},
	}

	httputil.JSONReponse(w, http.StatusCreated, payload)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req UserAuthRequest

	if !httputil.DecodeJSONBody[UserAuthRequest](w, r, &req) {
		return
	}

	if req.Email == "" || req.Password == "" {
		httputil.ErrorResponse(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := h.service.LoginUser(r.Context(), req.Email, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidLogin):
			httputil.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		default:
			httputil.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	atkToken, rtkToken, err := h.jwt.GenerateToken(user.ID)
	if err != nil {
		httputil.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	payload := UserResponse{
		Message: "User sign in successful",
		Auth: Auth{
			Atk: atkToken,
			Rtk: rtkToken,
		},
	}

	httputil.JSONReponse(w, http.StatusOK, payload)
}

func (h *AuthHandler) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}
