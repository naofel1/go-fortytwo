package fortytwo

import (
	"strconv"
	"time"
)

type CursusUserID int

func (pID CursusUserID) String() string {
	return strconv.Itoa(int(pID))
}

type CursusUsers []*CursusUser

type CursusUser struct {
	ID           int           `json:"id"`
	BeginAt      time.Time     `json:"begin_at"`
	EndAt        *time.Time    `json:"end_at,omitempty"`
	Grade        string        `json:"grade"`
	Level        float64       `json:"level"`
	Skills       []interface{} `json:"skills"`
	CursusId     int           `json:"cursus_id"`
	HasCoalition bool          `json:"has_coalition"`
	User         User          `json:"user"`
	Cursus       Cursus        `json:"cursus"`
}

type CursusUserQueryRequest struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}
