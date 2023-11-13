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
	tooBigError := fmt.Sprintf("Max file size is %d MB", app.Config.Limits.MaxFileSize>>20)

	// Parse the multipart form with a max file size
	err := r.ParseMultipartForm(app.Config.Limits.MaxFileSize)
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

	if fileHeader.Size > app.Config.Limits.MaxFileSize {
		http.Error(w, tooBigError, http.StatusBadRequest)
		return
	}

	metaFile, err := app.FileService.Upload(file, fileHeader)
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
	metaFile, err := app.FileService.Get(id)
	if err == service.ErrFileNotFound {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// Only return the safely exposable metadata
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
	metaFile, err := app.FileService.Get(id)
	if err == service.ErrFileNotFound {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// Retrieve the file
	file, err := app.FileService.LoadFile(metaFile)
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

// DeleteFile deletes the file with the given id from storage and the database.
func (app *App) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Check if the file exists and retrieve metadata
	metaFile, err := app.FileService.Get(id)
	if err == service.ErrFileNotFound {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// Delete the file
	err = app.FileService.Delete(metaFile)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	common.EncodeJSONAndSend(w, dto.DeleteFileResponse{
		ID: id,
	})
}

func (app *App) SearchFilesByDateRange(w http.ResponseWriter, r *http.Request) {
	fromParam := chi.URLParam(r, "from")
	toParam := chi.URLParam(r, "to")

	from, e1 := common.ParseAndValidateDate(fromParam)
	to, e2 := common.ParseAndValidateDate(toParam)

	// Validate the dates are in the correct format YYYY-MM-DD
	if e1 != nil || e2 != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	metaFiles, err := app.FileService.SearchByDateRange(from, to)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var res dto.SearchFilesResponse
	res.Count = len(metaFiles)
	res.Files = make([]dto.GetFileResponse, res.Count)
	for i, metaFile := range metaFiles {
		res.Files[i] = dto.GetFileResponse{
			ID:   metaFile.UUID,
			Name: metaFile.Name,
			Size: metaFile.Size,
		}
	}

	common.EncodeJSONAndSend(w, res)
}

func (app *App) DeleteFiles(w http.ResponseWriter, r *http.Request) {
	fromParam := chi.URLParam(r, "from")
	toParam := chi.URLParam(r, "to")

	from, e1 := common.ParseAndValidateDate(fromParam)
	to, e2 := common.ParseAndValidateDate(toParam)

	// Validate the dates are in the correct format YYYY-MM-DD
	if e1 != nil || e2 != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	metaFiles, err := app.FileService.SearchByDateRange(from, to)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var res dto.DeleteFilesResponse
	res.Count = len(metaFiles)
	res.Result = make([]dto.DeleteFilesResponseItem, res.Count)
	for i, metaFile := range metaFiles {
		res.Result[i] = dto.DeleteFilesResponseItem{
			ID: metaFile.UUID,
		}
		err = app.FileService.Delete(metaFile)
		if err != nil {
			// TODO: Better error handling here
			res.Result[i].Error = err.Error()
		}
	}

	common.EncodeJSONAndSend(w, res)
}
