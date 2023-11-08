package service

import (
	"dryve/internal/datastruct"
	"dryve/internal/repository"
)

type FileService interface {
	GetFile(id uint) (*datastruct.File, error)
}

type fileService struct {
	dao repository.DAO
}

func NewFileService(dao repository.DAO) FileService {
	return &fileService{dao: dao}
}

func (s *fileService) GetFile(id uint) (*datastruct.File, error) {
	return s.dao.NewFileQuery().GetFile(id)
}
