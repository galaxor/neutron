package memory

import (
	"errors"
	"time"

	"github.com/emersion/neutron/backend"
)

type SessionsBackend struct {
	sessions []*backend.Session
}

func (b *SessionsBackend) getSessionIndex(id string) (int, error) {
	for i, s := range b.sessions {
		if s.ID == id {
			return i, nil
		}
	}
	return -1, errors.New("No such session")
}

func (b *SessionsBackend) ListSessions(user string) (sessions []*backend.Session, err error) {
	for _, s := range b.sessions {
		if s.User.ID == user {
			sessions = append(sessions, s)
		}
	}
	return
}

func (b *SessionsBackend) GetSession(id string) (session *backend.Session, err error) {
	i, err := b.getSessionIndex(id)
	if err != nil {
		return
	}

	session = b.sessions[i]
	return
}

func (b *SessionsBackend) InsertSession(session *backend.Session) (*backend.Session, error) {
	session.ID = generateId()
	b.sessions = append(b.sessions, session)
	return session, nil
}

func (b *SessionsBackend) KeepSessionAlive(id string) error {
	session, err := b.GetSession(id)
	if err != nil {
		return err
	}

	session.Time = time.Now().Unix()
	return nil
}

func (b *SessionsBackend) DeleteSession(id string) error {
	i, err := b.getSessionIndex(id)
	if err != nil {
		return err
	}

	b.sessions = append(b.sessions[:i], b.sessions[i+1:]...)
	return nil
}

func NewSessionsBackend() backend.SessionsBackend {
	return &SessionsBackend{}
}
