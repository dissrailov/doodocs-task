package service

import (
	"archive/zip"
	"bytes"
	"doodocs-task/internal/models"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

var validMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func (s *service) AnalyzeArchive(file io.Reader, fileHeader *multipart.FileHeader) (models.ArchiveInfoResponse, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return models.ArchiveInfoResponse{}, fmt.Errorf("failed to read into memory: %v", err)
	}

	archive, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return models.ArchiveInfoResponse{}, fmt.Errorf("failed to open archive: %v", err)
	}

	var files []models.ArchiveFile
	var totalsize float64

	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			continue
		}

		fileInfo := models.ArchiveFile{
			FilePath: f.Name,
			Size:     float64(f.UncompressedSize64),
			MimeType: http.DetectContentType([]byte(f.Name)),
		}
		files = append(files, fileInfo)
		totalsize += float64(f.UncompressedSize64)
	}

	response := models.ArchiveInfoResponse{
		Filename:    fileHeader.Filename,
		ArchiveSize: float64(buf.Len()),
		TotalSize:   totalsize,
		TotalFiles:  float64(len(files)),
		Files:       files,
	}
	return response, nil
}

func (s *service) CreateArchive(files []*multipart.FileHeader) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		mimeType := getMimeType(file)
		if !validMimeTypes[mimeType] {
			return nil, fmt.Errorf("mime-type %s is not supported", mimeType)
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("failed to seek file: %v", err)
		}

		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create archive file: %v", err)
		}

		if _, err = io.Copy(zipFile, file); err != nil {
			return nil, fmt.Errorf("failed to copy file to archive: %v", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close archive writer: %v", err)
	}
	return buf.Bytes(), nil
}

func getMimeType(file multipart.File) string {
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return "application/octet-stream"
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "application/octet-stream"
	}
	return http.DetectContentType(buf)
}
