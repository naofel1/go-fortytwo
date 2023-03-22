package fortytwo

import (
	"time"
)

type SkillID int

type Skills []*Skill

type Skill struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
}

type SkillQueryRequest struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}
