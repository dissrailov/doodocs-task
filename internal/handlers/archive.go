package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *HandlerApp) AnalyzeArchive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file: %v", err), http.StatusBadRequest)
		return
	}

	defer file.Close()

	archiveInfo, err := h.ArchiveService.AnalyzeArchive(file, fileHeader)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to analyze archive: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(archiveInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *HandlerApp) CreateArchive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		http.Error(w, fmt.Sprintf("No files provided"), http.StatusBadRequest)
		return
	}

	archiveData, err := h.ArchiveService.CreateArchive(files)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create archive: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(archiveData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *HandlerApp) SendArchive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	emails := r.Form["emails"]
	if len(emails) == 0 {
		http.Error(w, fmt.Sprintf("No emails provided"), http.StatusBadRequest)
		return
	}

	err = h.ArchiveService.SendEmail("Subject: Your attachment", "Here is the file you requested.", emails, file, fileHeader)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
