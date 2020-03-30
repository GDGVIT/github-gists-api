package gists

type Gist struct {
	Url         string          `json:"url"`
	ID          string          `json:"id"`
	IsPublic    bool            `json:"public"`
	UpdatedAt   string          `json:"updated_at"`
	Description string          `json:"description"`
	Files       map[string]File `json:"files"`
}

type File struct {
	GistID      string `json:"gist_id"`
	GistUrl     string `json:"gist_url"`
	IsPublic    bool   `json:"public"`
	UpdatedAt   string `json:"updated_at"`
	Filename    string `json:"filename"`
	Language    string `json:"language"`
	RawUrl      string `json:"raw_url"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type GistFile struct {
	GistID      string `json:"gist_id"`
	Description string `json:"description"`
	IsPublic    bool   `json:"public"`
	Filename    string `json:"filename"`
	Content     string `json:"content"`
}

type CreateFileRequest struct {
	Description string                 `json:"description"`
	IsPublic    bool                   `json:"public"`
	Files       map[string]FileContent `json:"files"`
}

type FileContent struct {
	Content string `json:"content"`
}

type UpdateFileRequest struct {
	Description string                 `json:"description"`
	Files       map[string]FileContent `json:"files"`
}

type DeleteGist struct {
	GistID string `json:"gist_id"`
}
