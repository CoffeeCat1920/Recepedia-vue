package modals

import (
	"time"
)

type Session struct {
	SessionId string `json:"sessionid"`
	OwnerId   string `json:"ownerid"`
	Exp       string `json:"exp"`
}

func NewSession(OwnerId string) *Session {
	return &Session{
		SessionId: generateSessionToken(),
		OwnerId:   OwnerId,
		Exp:       time.Now().AddDate(0, 2, 0).Format("2006-01-02 15:04:05"),
	}
}

func (s *Session) IsExpired() bool {
	t, err := getTime(s.Exp)
	if err != nil {
		return false
	}
	return t.Before(time.Now())
}

func (s *Session) GetExpTime() (time.Time, error) {
	t, err := getTime(s.Exp)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}
