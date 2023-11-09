package app

import (
	"dryve/internal/app/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TODO: Move this stuff from here
var defaultFileStoragePath = "/tmp/hj-filestorage"

// GetFile handles the get file endpoint.
func (app *App) GetFile(w http.ResponseWriter, r *http.Request) {
	rid, _ := strconv.Atoi(chi.URLParam(r, "id"))
	id := uint(rid)

	w.WriteHeader(http.StatusCreated)
	common.EncodeJSONAndSend(w, map[string]any{"id": id})
}

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

	// Generate a UUID for the file
	// TODO: Validate file name against database to prevent duplicate filenames.
	//       e.g. Mechanism of write-to-reserve and commit-to-store.
	id := uuid.New().String()

	if fileHeader.Size > app.config.Limits.MaxFileSize {
		http.Error(w, tooBigError, http.StatusBadRequest)
		return
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(buff)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creates the uploads directory if it doesn't exist
	// TODO: Implement nested folders based on filename in a separate component
	//       to support large amounts of files on multiple locations/servers.
	//       e.g. 1234567890.jpg -> 123/456/7890.jpg
	err = os.MkdirAll(defaultFileStoragePath, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	storedFilename := fmt.Sprintf("%s%s", id, filepath.Ext(fileHeader.Filename))
	filePath := filepath.Join(defaultFileStoragePath, storedFilename)
	f, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a database entry for the file
	fileSize := fileHeader.Size
	uploadTime := time.Now().UTC()
	// TODO: Update DB
	fmt.Println(id, fileSize, uploadTime, filetype, storedFilename)

	// TODO: Add model for this response
	// Return a JSON response with the file ID
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"id":      id,
		"message": "File uploaded successfully",
	}
	json.NewEncoder(w).Encode(response)
}
