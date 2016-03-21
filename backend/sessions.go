package backend

// A session.
type Session struct {
	ID string
	Time int64
	User *User
}

// Stores sessions data to share them between different backends.
type SessionsBackend interface {
	// Lists all sessions of one user.
	ListSessions(user string) ([]*Session, error)
	// Get a session.
	GetSession(id string) (*Session, error)
	// Insert a new session.
	InsertSession(session *Session) (*Session, error)
	// Keeps a session alive.
	KeepSessionAlive(id string) error
	// Delete a session.
	DeleteSession(id string) error
}
