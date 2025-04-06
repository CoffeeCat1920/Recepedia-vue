package database

import (
	"big/internal/modals"
	"fmt"
)

func (s *service) AddAdminSession(session *modals.AdminSession) (error) {
  query := fmt.Sprintf(`
    INSERT INTO adminsessions(sessionid, exp) 
    VALUES('%s', '%s')
  `, session.SessionId, session.Exp)

  _, err := s.db.Exec(query) 

  if err != nil {
    return err
  }

  return nil
}

func (s *service) GetAdminSession(sessionId string) (*modals.AdminSession, error) {
  var session modals.AdminSession
  query := fmt.Sprintf("SELECT * FROM adminsessions WHERE sessionid = '%s';", sessionId)
  row := s.db.QueryRow(query)
  err := row.Scan(&session.SessionId, &session.Exp)

  if err != nil {
    return nil, err 
  }

  return &session, nil
}

func (s *service) DeleteAdminSession(sessionId string) error {
  _, err := s.db.Exec("DELETE FROM adminsessions WHERE sessionid = $1;", sessionId)
  if err != nil {
    fmt.Printf("Error deleting session: %v\n", err)
  }
  return err
}
