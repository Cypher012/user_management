package session

type SessionService struct {
	repo *SessionRepository
}

func NewSessionService(repo *SessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}
