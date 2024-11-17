package service

import (
	"doodocs-task/internal/models"
	"io"
	"mime/multipart"
)

type ArchiveServiceI interface {
	AnalyzeArchive(file io.Reader, fileHeader *multipart.FileHeader) (models.ArchiveInfoResponse, error)
	CreateArchive(files []*multipart.FileHeader) ([]byte, error)
}

type service struct{}

func NewService() ArchiveServiceI {
	return &service{}
}
