package gists

type Gist struct {
	Url       string            `json:"url"`
	ID        string            `json:"id"`
	IsPublic  bool              `json:"public"`
	UpdatedAt string            `json:"updated_at"`
	Files     []map[string]File `json:"files"`
}

type File struct {
	Filename string `json:"filename"`
	Language string `json:"language"`
	RawUrl   string `json:"raw_url"`
}
