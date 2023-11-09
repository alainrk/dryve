package service

import (
	"dryve/internal/datastruct"
	"dryve/internal/repository"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// TODO: Move this stuff from here
var defaultFileStoragePath = "/tmp/hj-filestorage"

var FileBadRequestError = fmt.Errorf("bad file request")
var FileProcessingError = fmt.Errorf("file processing error")

type FileService interface {
	Upload(multipart.File, *multipart.FileHeader) (datastruct.File, error)
}

type fileService struct {
	dao repository.DAO
}

func NewFileService(dao repository.DAO) FileService {
	return &fileService{dao: dao}
}

func (s *fileService) Upload(file multipart.File, fileHeader *multipart.FileHeader) (datastruct.File, error) {
	var metaFile datastruct.File

	// Generate a UUID for the file
	// TODO: Validate file name against database to prevent duplicate filenames.
	//       e.g. Mechanism of write-to-reserve and commit-to-store.
	id := uuid.New().String()

	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		return metaFile, FileProcessingError
	}

	filetype := http.DetectContentType(buff)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return metaFile, FileProcessingError
	}

	// Creates the uploads directory if it doesn't exist
	// TODO: Implement nested folders based on filename in a separate component
	//       to support large amounts of files on multiple locations/servers.
	//       e.g. 1234567890.jpg -> 123/456/7890.jpg
	err = os.MkdirAll(defaultFileStoragePath, os.ModePerm)
	if err != nil {
		return metaFile, FileProcessingError
	}
	storedFilename := fmt.Sprintf("%s%s", id, filepath.Ext(fileHeader.Filename))
	filePath := filepath.Join(defaultFileStoragePath, storedFilename)
	f, err := os.Create(filePath)
	if err != nil {
		return metaFile, FileBadRequestError
	}

	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return metaFile, FileProcessingError
	}

	// Create a database entry for the file
	fileSize := fileHeader.Size
	uploadTime := time.Now().UTC()
	// TODO: Update DB
	fmt.Println(id, fileSize, uploadTime, filetype, storedFilename)

	return metaFile, nil
	// return s.dao.NewFileQuery().Upload(id)
}
