package app

import (
	"dryve/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// UploadFile handles the upload file endpoint.
// It parses the multipart form, saves the file to disk, and creates a database entry for the file.
func (app *App) UploadFile(w http.ResponseWriter, r *http.Request) {
	tooBigError := fmt.Sprintf("Max file size is %d MB", app.config.Limits.MaxFileSize>>20)

	// Parse the multipart form with a max file size
	err := r.ParseMultipartForm(app.config.Limits.MaxFileSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s\n%s", err.Error(), tooBigError), http.StatusBadRequest)
		return
	}

	// Retrieve the file from the multipart form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if fileHeader.Size > app.config.Limits.MaxFileSize {
		http.Error(w, tooBigError, http.StatusBadRequest)
		return
	}

	metFile, err := app.fileService.Upload(file, fileHeader)
	if err == service.FileBadRequestError {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err == service.FileProcessingError {
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	// TODO: Add model for this response
	// Return a JSON response with the file ID
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"id":      metFile.ID,
		"message": "File uploaded successfully",
	}
	json.NewEncoder(w).Encode(response)
}
