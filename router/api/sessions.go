package api

import (
	"time"
)

const SessionTimeout = 10 * time.Minute

type Session struct {
	ID string
	UserID string
	Token string
	Timeout *time.Timer
}

func NewSession(user string, expire func()) *Session {
	return &Session{
		ID: generateId(),
		Token: generateId(),
		UserID: user,
		Timeout: time.AfterFunc(SessionTimeout, expire),
	}
}
