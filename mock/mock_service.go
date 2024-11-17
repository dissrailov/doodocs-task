package mock

import (
	"doodocs-task/internal/models"
	"io"
	"mime/multipart"
	"testing"
)

type MockService struct {
	AnalyzeArchiveFunc func(file io.Reader, fileHeader *multipart.FileHeader) (models.ArchiveInfoResponse, error)
	CreateArchiveFunc  func(files []*multipart.FileHeader) ([]byte, error)
	SendEmailFunc      func(subject string, body string, to []string, file multipart.File, fileHeader *multipart.FileHeader) error
}

func NewMockService(t *testing.T) *MockService {
	return &MockService{}
}

func (m *MockService) AnalyzeArchive(file io.Reader, header *multipart.FileHeader) (models.ArchiveInfoResponse, error) {
	if m.AnalyzeArchiveFunc == nil {
		return m.AnalyzeArchiveFunc(file, header)
	}
	return models.ArchiveInfoResponse{
		Filename:    header.Filename,
		ArchiveSize: float64(header.Size),
		TotalSize:   float64(header.Size),
		TotalFiles:  1,
		Files: []models.ArchiveFile{
			{
				FilePath: header.Filename,
				Size:     float64(header.Size),
				MimeType: "application/zip",
			},
		},
	}, nil
}

func (m *MockService) CreateArchive(files []*multipart.FileHeader) ([]byte, error) {
	if m.CreateArchiveFunc == nil {
		return m.CreateArchiveFunc(files)
	}
	return []byte("data"), nil
}

func (m *MockService) SendEmail(subject string, body string, to []string, file multipart.File, fileHeader *multipart.FileHeader) error {
	if m.SendEmailFunc == nil {
		return m.SendEmailFunc(subject, body, to, file, fileHeader)
	}
	return nil
}
