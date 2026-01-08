package auth

import "context"

type AuthService struct {
	repo *AuthRepository
}

func NewAuthService(repo *AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, email, password string) {

}
