package app

import (
	"dryve/internal/app/common"
	"dryve/internal/dto"
	"dryve/internal/service"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	metaFile, err := app.fileService.Upload(file, fileHeader)
	if err == service.ErrFileBadRequest {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err == service.ErrFileProcessing {
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	common.EncodeJSONAndSend(w, dto.UploadFileResponse{
		ID: metaFile.UUID,
	})
}

// GetFile returns the file with the given id (internal UUID).
func (app *App) GetFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Check if the file exists and retrieve metadata
	metaFile, err := app.fileService.Get(id)
	if err == service.ErrFileNotFound {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	common.EncodeJSONAndSend(w, dto.GetFileResponse{
		ID:   metaFile.UUID,
		Name: metaFile.Name,
		Size: metaFile.Size,
	})
}

// DownloadFile returns the file with the given id (internal UUID).
func (app *App) DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Check if the file exists and retrieve metadata
	metaFile, err := app.fileService.Get(id)
	if err == service.ErrFileNotFound {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// Retrieve the file
	file, err := app.fileService.LoadFile(metaFile)
	if err != nil {
		http.Error(w, "Internal error loading file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", metaFile.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", metaFile.Size))

	// Copy the file to the response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
