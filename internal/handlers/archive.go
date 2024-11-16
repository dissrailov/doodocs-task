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
