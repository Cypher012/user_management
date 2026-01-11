package authhttp

import (
	"errors"
	"log"
	"net/http"

	"github.com/Cypher012/userauth/internal/auth"
	"github.com/Cypher012/userauth/internal/email"
	"github.com/Cypher012/userauth/internal/http/httputil"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	service *auth.AuthService
	email   *email.EmailService
	jwt     *auth.JWTAuth
}

func NewAuthHandler(service *auth.AuthService, email *email.EmailService, jwt *auth.JWTAuth) *AuthHandler {
	return &AuthHandler{
		service: service,
		email:   email,
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

	rawToken, err := h.service.CreateEmailVerificationToken(
		r.Context(),
		user.ID,
	)
	if err != nil {
		httputil.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	auth.SetRefreshCookies(w, rtkToken)

	go func(email, token string) {
		log.Printf("sending verify email to %s...", email)
		if err := h.email.SendVerifyEmail(email, token); err != nil {
			log.Printf("verify email failed: %v", err)
		} else {
			log.Printf("verify email sent successfully to %s", email)
		}
	}(user.Email, rawToken)

	payload := UserResponse{
		Message: "User sign up succesful, check mail to verify email",
		Auth: Auth{
			Token: atkToken,
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

	auth.SetRefreshCookies(w, rtkToken)

	payload := UserResponse{
		Message: "User sign in successful",
		Auth: Auth{
			Token: atkToken,
		},
	}

	httputil.JSONReponse(w, http.StatusOK, payload)
}

func (h *AuthHandler) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rawToken := chi.URLParam(r, "token")

	if err := h.service.VerifyEmailVerificationToken(r.Context(), rawToken); err != nil {
		http.ServeFile(w, r, "internal/web/verify_error.html")
		return
	}

	http.ServeFile(w, r, "internal/web/verify_success.html")
}

func (h *AuthHandler) ResendVerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	claims, err := h.jwt.FromContext(r.Context())
	if err != nil {
		httputil.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
	}

	userId := claims.UserId

	rawToken, email, err := h.service.CreateResendEmailVerificationToken(r.Context(), userId)

	if err != nil {
		httputil.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	go func(email, token string) {
		log.Printf("sending verify email to %s...", email)
		if err := h.email.SendVerifyEmail(email, token); err != nil {
			log.Printf("verify email failed: %v", err)
		} else {
			log.Printf("verify email sent successfully to %s", email)
		}
	}(email, rawToken)

	payload := map[string]string{
		"Message": "Check email to veify email",
	}

	httputil.JSONReponse(w, http.StatusOK, payload)
}

func (h *AuthHandler) ForgetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req ForgetPasswordRequest

	if !httputil.DecodeJSONBody[ForgetPasswordRequest](w, r, &req) {
		return
	}

	if req.Email == "" {
		httputil.ErrorResponse(w, http.StatusBadRequest, "Email field is required")
		return
	}

	rawToken, err := h.service.CreateForgetPasswordToken(r.Context(), req.Email)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			httputil.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		default:
			httputil.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	go func(email, token string) {
		log.Printf("sending forget password email to %s...", email)
		if err := h.email.SendForgetPasswordEmail(email, token); err != nil {
			log.Printf("forget password email failed: %v", err)
		} else {
			log.Printf("forget password email sent successfully to %s", email)
		}
	}(req.Email, rawToken)

	payload := map[string]string{
		"Message": "Check email to reset password",
	}

	httputil.JSONReponse(w, http.StatusOK, payload)
}

func (h *AuthHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rawToken := chi.URLParam(r, "token")
	var req ResetPasswordRequest

	if !httputil.DecodeJSONBody[ResetPasswordRequest](w, r, &req) {
		return
	}

	userId, err := h.service.VerifyResetPasswordToken(r.Context(), rawToken)
	if err != nil {
		httputil.ErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	if err := h.service.ChangePassword(r.Context(), userId, req.Password); err != nil {
		switch {
		case errors.Is(err, auth.ErrPasswordHash):
			httputil.ErrorResponse(w, http.StatusFailedDependency, err.Error())
		default:
			httputil.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	payload := map[string]string{
		"Message": "Password has been changed successfully",
	}

	httputil.JSONReponse(w, http.StatusOK, payload)
}

func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}
