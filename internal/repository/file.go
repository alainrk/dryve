package repository

import (
	"dryve/internal/datastruct"

	"gorm.io/gorm"
)

type FileQuery interface {
	GetFile(id uint) (*datastruct.File, error)
}

type fileQuery struct {
	db *gorm.DB
}

func (d *dao) NewFileQuery() FileQuery {
	return &fileQuery{d.db}
}

func (u *fileQuery) GetFile(id uint) (*datastruct.File, error) {
	var file datastruct.File
	err := u.db.First(&file, id).Error
	return &file, err
}
