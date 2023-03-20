package fortytwo

import (
	"strconv"
)

type TitleID int

func (pID TitleID) String() string {
	return strconv.Itoa(int(pID))
}

type Titles []*Title

type Title struct {
	ID   TitleID `json:"id"`
	Name string  `json:"name"`
}
