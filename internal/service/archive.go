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
