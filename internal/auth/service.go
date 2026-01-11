package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Cypher012/userauth/internal/security"
	"github.com/Cypher012/userauth/internal/token"
)

var (
	ErrPasswordHash = errors.New("password hashing failed")
	ErrInvalidLogin = errors.New("invalid email or password")
)

type User struct {
	ID         string
	Email      string
	IsVerified bool
	IsActive   bool
	CreatedAt  time.Time
}

type AuthService struct {
	repo                    *AuthRepository
	token                   *token.TokenService
	verifyEmailTokenType    token.TokenType
	forgetPasswordTokenType token.TokenType
}

func NewAuthService(repo *AuthRepository, tokenSvc *token.TokenService) *AuthService {
	return &AuthService{
		repo:                    repo,
		token:                   tokenSvc,
		verifyEmailTokenType:    token.VerifyEmailTokenType,
		forgetPasswordTokenType: token.ForgetPasswordTokenType,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, email, password string) (User, error) {
	_, err := s.repo.GetUserByEmail(ctx, email)
	switch {
	case err == nil:
		return User{}, ErrUserAlreadyExists
	case !errors.Is(err, ErrUserNotFound):
		log.Println(err.Error())
		return User{}, err
	}

	hashedPassword, err := security.GenerateHashPassword(password)
	if err != nil {
		return User{}, ErrPasswordHash
	}

	createUserRow, err := s.repo.CreateUser(ctx, email, hashedPassword)

	if err != nil {
		return User{}, err
	}

	return User{
		ID:         createUserRow.ID.String(),
		Email:      createUserRow.Email,
		IsVerified: createUserRow.IsVerified,
		IsActive:   createUserRow.IsActive,
		CreatedAt:  createUserRow.CreatedAt.Time,
	}, nil
}

func (s *AuthService) LoginUser(ctx context.Context, email, password string) (User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, ErrInvalidLogin
	}

	if err := security.ComparePassword(user.PasswordHash, password); err != nil {
		return User{}, ErrInvalidLogin
	}

	return User{
		ID:         user.ID.String(),
		Email:      user.Email,
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt.Time,
	}, nil
}

func (s *AuthService) CreateEmailVerificationToken(
	ctx context.Context,
	userID string,
) (string, error) {
	return s.token.CreateToken(ctx, userID, s.verifyEmailTokenType)
}

func (s *AuthService) VerifyEmailVerificationToken(ctx context.Context, rawToken string) error {
	userId, err := s.token.VerifyToken(ctx, rawToken, s.verifyEmailTokenType)
	if err != nil {
		return err
	}
	return s.repo.SetUserEmailVerified(ctx, userId)
}

func (s *AuthService) CreateForgetPasswordToken(
	ctx context.Context,
	email string,
) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrUserNotFound
	}

	return s.token.CreateToken(ctx, user.ID.String(), s.forgetPasswordTokenType)
}

func (s *AuthService) VerifyResetPasswordToken(ctx context.Context, rawToken string) (string, error) {
	userID, err := s.token.VerifyToken(ctx, rawToken, s.forgetPasswordTokenType)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userId, password string) error {
	hashed, err := security.GenerateHashPassword(password)
	if err != nil {
		return ErrPasswordHash
	}
	return s.repo.UpdateUserPassword(ctx, userId, hashed)
}

func (s *AuthService) CreateResendEmailVerificationToken(
	ctx context.Context,
	userID string,
) (rawToken string, email string, err error) {
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return "", "", ErrUserNotFound
	}
	rawToken, err = s.token.CreateToken(ctx, userID, s.verifyEmailTokenType)
	if err != nil {
		return "", "", err
	}

	return rawToken, user.Email, nil
}
