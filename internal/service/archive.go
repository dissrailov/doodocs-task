package service

import (
	"archive/zip"
	"bytes"
	"doodocs-task/internal/models"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"path/filepath"
	"strings"
)

var validMimeTypesForCreate = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}
var validMimeTypesForSend = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/pdf": true,
}

func (s *service) AnalyzeArchive(file io.Reader, fileHeader *multipart.FileHeader) (models.ArchiveInfoResponse, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return models.ArchiveInfoResponse{}, fmt.Errorf("failed to read into memory: %v", err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	if mimeType != "application/zip" {
		return models.ArchiveInfoResponse{}, fmt.Errorf("unsupported file type: %s", mimeType)
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

		mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
		if !validMimeTypesForCreate[mimeType] {
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

func (s *service) SendEmail(subject string, body string, to []string, file multipart.File, fileHeader *multipart.FileHeader) error {
	fileContent, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer fileContent.Close()

	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	if !validMimeTypesForSend[mimeType] {
		return fmt.Errorf("mime-type %s is not supported", mimeType)
	}

	fileData, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return fmt.Errorf("failed to read file content: %v", err)
	}

	message, err := s.buildEmailMessage(subject, body, to, fileHeader.Filename, mimeType, fileData)
	if err != nil {
		return fmt.Errorf("failed to build email message: %v", err)
	}
	if err := s.sendSMTPMessage(to, message); err != nil {
		return fmt.Errorf("failed to send email message: %v", err)
	}

	return nil
}

func (s *service) buildEmailMessage(subject, body string, to []string, filename, mimeType string, fileData []byte) (string, error) {
	boundary := "boundary-12345"

	message := fmt.Sprintf("From: %s\r\n", s.SMTP.From)
	message += fmt.Sprintf("To: %s\r\n", strings.Join(to, ","))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary)
	message += "\r\n--" + boundary + "\r\n"
	message += fmt.Sprintf("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	message += "Content-Transfer-Encoding: 7bit\r\n"
	message += "\r\n" + body + "\r\n"
	message += "\r\n--" + boundary + "\r\n"

	attachment := fmt.Sprintf("Content-Type: %s; name=%s\r\n", mimeType, filename)
	attachment += fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", filename)
	attachment += "Content-Transfer-Encoding: base64\r\n"
	attachment += "\r\n"
	attachment += base64.StdEncoding.EncodeToString(fileData)
	attachment += "\r\n--" + boundary + "--\r\n"

	message += attachment

	return message, nil
}

func (s *service) sendSMTPMessage(to []string, message string) error {
	auth := smtp.PlainAuth("", s.SMTP.User, s.SMTP.Password, s.SMTP.Host)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", s.SMTP.Host, s.SMTP.Port), auth, s.SMTP.From, to, []byte(message))
	return err
}
