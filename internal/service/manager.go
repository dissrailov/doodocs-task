package service

import (
	"doodocs-task/internal/models"
	"io"
	"mime/multipart"
)

type ArchiveServiceI interface {
	AnalyzeArchive(file io.Reader, fileHeader *multipart.FileHeader) (models.ArchiveInfoResponse, error)
}

type service struct{}

func NewService() ArchiveServiceI {
	return &service{}
}
