package database

import (
	"big/internal/modals"
)

func (s* service) AddSession(session *modals.Session) (error) {
  query := `INSERT INTO sessions(sessionid, ownerid, exp) 
    VALUES($1, $2, $3)`

  _, err := s.db.Exec(query, session.SessionId, session.OwnerId, session.Exp) 

  if err != nil {
    return err
  }

  return nil
}

func (s *service)DeleteSession(sessionId string) (error) {
  query := `DELETE FROM sessions WHERE sessionId = $1`

  _, err := s.db.Exec(query, sessionId)

  if err != nil {
    return err
  }

  return nil
}

func (s *service) DeleteSessionByUser(ownerId string) (error) {
	query := `DELETE FROM sessions WHERE ownerid = $1`

  _, err := s.db.Exec(query, ownerId)

  if err != nil {
    return err
  }

  return nil
}   

func (s *service)GetSession(sessionId string) (*modals.Session, error) {
  var session modals.Session
  query := "SELECT * FROM sessions WHERE sessionid = $1;"
  row := s.db.QueryRow(query, sessionId)
  err := row.Scan(&session.SessionId, &session.OwnerId, &session.Exp)

  if err != nil {
    return nil, err 
  }

  return &session, nil
}
