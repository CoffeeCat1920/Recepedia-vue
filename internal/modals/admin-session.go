package modals

import "time"

type AdminSession struct {
	SessionId string `json:"sessionId"`
	Exp       string `json:"exp"`
}

func NewAdminSession() *AdminSession {
	return &AdminSession{
		SessionId: generateSessionToken(),
		Exp:       time.Now().AddDate(0, 2, 0).Format("2006-01-02 15:04:05"),
	}
}

func (s *AdminSession) IsExpired() bool {
	t, err := getTime(s.Exp)
	if err != nil {
		return false
	}
	return t.Before(time.Now())
}

func (s *AdminSession) GetExpTime() (time.Time, error) {
	t, err := getTime(s.Exp)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}
