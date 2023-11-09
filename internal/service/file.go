package service

import (
	"dryve/internal/datastruct"
	"dryve/internal/repository"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: Move this stuff from here
var defaultFileStoragePath = "/tmp/hj-filestorage"

var ErrFileNotFound = fmt.Errorf("file not found")
var ErrFileBadRequest = fmt.Errorf("bad file request")
var ErrFileProcessing = fmt.Errorf("file processing error")
var ErrFileInternal = fmt.Errorf("file processing error")

type FileService interface {
	Get(id string) (*datastruct.File, error)
	Upload(multipart.File, *multipart.FileHeader) (*datastruct.File, error)
}

type fileService struct {
	dao repository.DAO
}

func NewFileService(dao repository.DAO) FileService {
	return &fileService{dao: dao}
}

func (s *fileService) Get(id string) (*datastruct.File, error) {
	var metaFile *datastruct.File

	metaFile, err := s.dao.NewFileQuery().Get(id)
	// TODO: Remove this dependency for an internal error instead
	if err == gorm.ErrRecordNotFound {
		return metaFile, ErrFileNotFound
	}
	if err != nil {
		return metaFile, ErrFileInternal
	}

	return metaFile, nil
}

func (s *fileService) Upload(file multipart.File, fileHeader *multipart.FileHeader) (*datastruct.File, error) {
	var metaFile *datastruct.File

	// Generate a UUID for the file
	// TODO: Validate file name against database to prevent duplicate filenames.
	//       e.g. Mechanism of write-to-reserve and commit-to-store.
	id := uuid.New().String()

	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		return metaFile, ErrFileProcessing
	}

	// TODO: Validate/Restrict available file type
	// filetype := http.DetectContentType(buff)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return metaFile, ErrFileProcessing
	}

	// Creates the uploads directory if it doesn't exist
	// TODO: Implement nested folders based on filename in a separate component
	//       to support large amounts of files on multiple locations/servers.
	//       e.g. 1234567890.jpg -> 123/456/7890.jpg
	err = os.MkdirAll(defaultFileStoragePath, os.ModePerm)
	if err != nil {
		return metaFile, ErrFileProcessing
	}
	storedFilename := fmt.Sprintf("%s%s", id, filepath.Ext(fileHeader.Filename))
	filePath := filepath.Join(defaultFileStoragePath, storedFilename)
	f, err := os.Create(filePath)
	if err != nil {
		return metaFile, ErrFileBadRequest
	}

	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return metaFile, ErrFileProcessing
	}

	// Create a database entry for the file
	fileSize := fileHeader.Size

	metaFile, err = s.dao.NewFileQuery().Create(id, fileHeader.Filename, fileSize, storedFilename)
	if err != nil {
		return metaFile, ErrFileProcessing
	}

	return metaFile, nil
}
