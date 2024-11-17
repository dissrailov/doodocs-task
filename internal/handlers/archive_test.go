package handlers

import (
	"bytes"
	"doodocs-task/internal/models"
	"doodocs-task/mock"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAnalyzeArchive(t *testing.T) {
	mockService := mock.NewMockService(t)
	mockService.AnalyzeArchiveFunc = func(file io.Reader, fileHeader *multipart.FileHeader) (models.ArchiveInfoResponse, error) {
		return models.ArchiveInfoResponse{
			Filename:    fileHeader.Filename,
			ArchiveSize: float64(fileHeader.Size),
			TotalSize:   float64(fileHeader.Size),
			TotalFiles:  1,
			Files: []models.ArchiveFile{
				{
					FilePath: fileHeader.Filename,
					Size:     float64(fileHeader.Size),
					MimeType: "application/zip",
				},
			},
		}, nil
	}
	handler := &HandlerApp{
		ArchiveService: mockService,
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "test.zip")
	part.Write([]byte("dummy data"))
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/archive/information", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.AnalyzeArchive(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	expectedResponse := `{"filename":"test.zip","archiveSize":10,"totalSize":10,"totalFiles":1,"files":[{"filePath":"test.zip","size":10,"mimeType":"application/zip"}]}`
	if strings.TrimSpace(rr.Body.String()) != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/zip")
	}
}
func TestCreateArchive(t *testing.T) {
	mockService := mock.NewMockService(t)
	mockService.CreateArchiveFunc = func(files []*multipart.FileHeader) ([]byte, error) {
		return []byte{0x50, 0x4b, 0x03, 0x04, 0x14, 0x00, 0x00, 0x00}, nil
	}
	handler := &HandlerApp{
		ArchiveService: mockService,
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("files[]", "test.jpg")
	if err != nil {
		t.Fatalf("could not create form file: %v", err)
	}
	part.Write([]byte("dummy data"))

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/archive/files", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.CreateArchive(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/zip" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/zip")
	}
}
func TestSendEmail(t *testing.T) {
	mockService := mock.NewMockService(t)
	mockService.SendEmailFunc = func(subject string, body string, to []string, file multipart.File, fileHeader *multipart.FileHeader) error {
		if len(to) == 0 {
			return fmt.Errorf("no to")
		}
		return nil
	}
	handler := &HandlerApp{
		ArchiveService: mockService,
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "test.pdf")
	part.Write([]byte("dummy data"))
	writer.WriteField("emails", "test1@example.com")
	writer.WriteField("emails", "test2@example.com")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/mail/file", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler.SendArchive(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}
