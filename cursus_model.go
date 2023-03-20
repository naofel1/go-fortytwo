package fortytwo

import (
	"time"
)

type CursusID int

type CursusSlice []*Cursus

type Cursus struct {
	ID        int        `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Kind      CursusKind `json:"kind"`
}

type CursusQueryRequest struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}
