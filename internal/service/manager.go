package service

import (
	"doodocs-task/config"
	"doodocs-task/internal/models"
	"io"
	"mime/multipart"
)

type ArchiveServiceI interface {
	AnalyzeArchive(file io.Reader, fileHeader *multipart.FileHeader) (models.ArchiveInfoResponse, error)
	CreateArchive(files []*multipart.FileHeader) ([]byte, error)
	SendEmail(subject string, body string, to []string, file multipart.File, fileHeader *multipart.FileHeader) error
}

type service struct {
	SMTP config.SMTP
}

func NewService(smtp *config.SMTP) ArchiveServiceI {
	return &service{
		SMTP: *smtp,
	}
}
