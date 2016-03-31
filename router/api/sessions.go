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
		ID: "session_id", // TODO: generate this
		Token: "access_token", // TODO: generate this
		UserID: user,
		Timeout: time.AfterFunc(SessionTimeout, expire),
	}
}
