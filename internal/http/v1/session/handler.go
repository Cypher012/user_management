package session

import (
	"net/http"

	"github.com/Cypher012/user_management/internal/session"
)

type SessionHandler struct {
	service *session.SessionService
}

func NewSessionHandler(service *session.SessionService) *SessionHandler {
	return &SessionHandler{
		service: service,
	}
}

func (h *SessionHandler) ListAllActiveSessionsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

}
