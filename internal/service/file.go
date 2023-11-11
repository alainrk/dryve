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

var ErrFileNotFound = fmt.Errorf("file not found")
var ErrFileBadRequest = fmt.Errorf("bad file request")
var ErrFileProcessing = fmt.Errorf("file processing error")
var ErrFileInternal = fmt.Errorf("file processing error")

type FileService interface {
	Get(id string) (*datastruct.File, error)
	Upload(multipart.File, *multipart.FileHeader) (*datastruct.File, error)
	Delete(metaFile *datastruct.File) error
	LoadFile(metaFile *datastruct.File) (io.ReadCloser, error)
}

type fileService struct {
	dao             repository.DAO
	fileStoragePath string
}

func NewFileService(dao repository.DAO, path string) FileService {
	return &fileService{
		dao:             dao,
		fileStoragePath: path,
	}
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
	err = os.MkdirAll(s.fileStoragePath, os.ModePerm)
	if err != nil {
		return metaFile, ErrFileProcessing
	}
	storedFilename := fmt.Sprintf("%s%s", id, filepath.Ext(fileHeader.Filename))
	filePath := filepath.Join(s.fileStoragePath, storedFilename)
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

func (s *fileService) Delete(metaFile *datastruct.File) error {
	filePath := filepath.Join(s.fileStoragePath, metaFile.Filename)
	err := os.Remove(filePath)
	if err != nil {
		return ErrFileInternal
	}

	// Remove from the database through dto
	err = s.dao.NewFileQuery().Delete(metaFile.UUID)
	if err != nil {
		return ErrFileInternal
	}

	return nil
}

func (s *fileService) LoadFile(metaFile *datastruct.File) (file io.ReadCloser, err error) {
	filePath := filepath.Join(s.fileStoragePath, metaFile.Filename)
	file, err = os.Open(filePath)
	if err != nil {
		// TODO: Better management of different errors
		return nil, ErrFileInternal
	}

	return file, nil
}
