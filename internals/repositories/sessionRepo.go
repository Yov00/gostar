package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"templ_workout/internals/models"
)

type SessionRepo struct {
	DB *sql.DB
}

func (s *SessionRepo) Insert(session models.Session) error {
	var err error
	_, err = s.DB.Exec(
		"INSERT INTO sessions(user_id,session_token,csrf_token,user_agent,ip_address,created_at,expires_at) VALUES(?,?,?,?,?,?,?)",
		session.UserID, session.SessionToken, session.CSRFToken, session.UserAgent, session.IPAddress, session.CreatedAt, session.ExpiresAt)

	return err
}

func (s *SessionRepo) GetSessionByCSRFAndSessionToken(session_token string, csrf string) (*models.Session, error) {
	var err error
	var session models.Session

	row := s.DB.QueryRow("SELECT id,user_id, session_token, csrf_token, ip_address,user_agent, created_at,expires_at  from sessions where session_token = ? and csrf_token = ? AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP) ", strings.TrimSpace(session_token), strings.TrimSpace(csrf))

	err = row.Scan(&session.ID, &session.UserID, &session.SessionToken, &session.CSRFToken, &session.IPAddress, &session.UserAgent, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		fmt.Printf("Could not find user session %v", err)
		return nil, err
	}

	return &session, nil
}
