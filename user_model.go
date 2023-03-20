package fortytwo

import (
	"time"
)

type UserID int

type LocationsStat interface{}

type Users []User

type User struct {
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	ID              int           `json:"id"`
	Email           string        `json:"email"`
	Login           string        `json:"login"`
	FirstName       string        `json:"first_name"`
	LastName        string        `json:"last_name"`
	UsualFullName   string        `json:"usual_full_name"`
	UsualFirstName  string        `json:"usual_first_name"`
	Url             string        `json:"url"`
	Phone           string        `json:"phone"`
	Displayname     string        `json:"displayname"`
	Kind            string        `json:"kind"`
	Image           Image         `json:"image"`
	Staff           bool          `json:"staff?"`
	CorrectionPoint int           `json:"correction_point"`
	PoolMonth       string        `json:"pool_month"`
	PoolYear        string        `json:"pool_year"`
	Location        interface{}   `json:"location"`
	Wallet          int           `json:"wallet"`
	AnonymizeDate   time.Time     `json:"anonymize_date"`
	DataErasureDate time.Time     `json:"data_erasure_date,omitempty"`
	AlumnizedAt     time.Time     `json:"alumnized_at,omitempty"`
	Alumni          bool          `json:"alumni?"`
	Active          bool          `json:"active?"`
	Groups          []interface{} `json:"groups"`
	CursusUsers     []CursusUser  `json:"cursus_users"`
	ProjectsUsers   []interface{} `json:"projects_users"`
	LanguagesUsers  []struct {
		ID         int       `json:"id"`
		LanguageID int       `json:"language_id"`
		UserID     int       `json:"user_id"`
		Position   int       `json:"position"`
		CreatedAt  time.Time `json:"created_at"`
	}
}
