package fortytwo

type ProjectID int

type Projects []*Project

type Project struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Difficulty  int           `json:"difficulty"`
	Description string        `json:"description"`
	Parent      interface{}   `json:"parent"`
	Children    []interface{} `json:"children"`
	Objectives  []string      `json:"objectives"`
	Attachments []interface{} `json:"attachments"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
	Exam        bool          `json:"exam"`
}

type ProjectQueryRequest struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}
