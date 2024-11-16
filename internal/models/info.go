package models

type ArchiveInfoResponse struct {
	Filename    string        `json:"filename"`
	ArchiveSize float64       `json:"archiveSize"`
	TotalSize   float64       `json:"totalSize"`
	TotalFiles  float64       `json:"totalFiles"`
	Files       []ArchiveFile `json:"files"`
}

type ArchiveFile struct {
	FilePath string  `json:"filePath"`
	Size     float64 `json:"size"`
	MimeType string  `json:"mimeType"`
}
