package fortytwo

type ImageVersion struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
	Micro  string `json:"micro"`
}

type Image struct {
	Link     string        `json:"link"`
	Versions *ImageVersion `json:"versions"`
}
