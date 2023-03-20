package fortytwo

type Achievements []Achievement

type Achievement struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Tier         string   `json:"tier"`
	Kind         string   `json:"kind"`
	Visible      bool     `json:"visible"`
	Image        string   `json:"image"`
	NbrOfSuccess int      `json:"nbr_of_success,omitempty"`
	UsersURL     string   `json:"users_url,omitempty"`
	Achievements []string `json:"achievements,omitempty"`
	Campus       []string `json:"campus,omitempty"`

	Parent *Achievement `json:"parent,omitempty"`
	Title  *Title       `json:"title,omitempty"`
}

type AchievementQueryRequest struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}
