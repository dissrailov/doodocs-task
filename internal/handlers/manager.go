package handlers

import "doodocs-task/internal/service"

type HandlerApp struct {
	ArchiveService service.ArchiveServiceI
}

func NewHandler(archiveService service.ArchiveServiceI) *HandlerApp {
	return &HandlerApp{
		ArchiveService: archiveService,
	}
}
