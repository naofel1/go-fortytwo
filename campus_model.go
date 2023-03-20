package fortytwo

import (
	"strconv"
)

type CampusID int

func (pID CampusID) String() string {
	return strconv.Itoa(int(pID))
}

type CampusSlice []*Campus

type Campus struct {
	ID   CampusID `json:"id"`
	Name string   `json:"name"`
}
