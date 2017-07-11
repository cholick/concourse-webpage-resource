package models

type Source struct {
	Url      string `json:"url"`
	Filename string `json:"filename"`
}

type Version struct {
	ETag         string `json:"eTag,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version Version `json:"version"`
}
