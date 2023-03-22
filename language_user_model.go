package fortytwo

import "time"

type LanguageUserID int

type LanguageUsers []*LanguageUser

type LanguageUser struct {
	ID         int       `json:"id"`
	LanguageID int       `json:"language_id"`
	UserID     int       `json:"user_id"`
	Position   int       `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
}
