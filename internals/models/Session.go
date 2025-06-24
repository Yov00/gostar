package models

import "time"

type Session struct {
	ID           int        `db:"id" json:"id"`
	UserID       int        `db:"user_id" json:"user_id"`
	SessionToken string     `db:"session_token" json:"session_token"`
	CSRFToken    string     `db:"csrf_token" json:"csrf_token"`
	UserAgent    *string    `db:"user_agent" json:"user_agent,omitempty"`
	IPAddress    *string    `db:"ip_address" json:"ip_address,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	ExpiresAt    *time.Time `db:"expires_at" json:"expires_at,omitempty"`
}

