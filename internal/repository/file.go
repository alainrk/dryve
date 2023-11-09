package repository

import (
	"dryve/internal/datastruct"

	"gorm.io/gorm"
)

type FileQuery interface {
	Create(UUID string, Name string, Size int64, Filename string) (*datastruct.File, error)
	Get(id uint) (*datastruct.File, error)
}

type fileQuery struct {
	db *gorm.DB
}

func (d *dao) NewFileQuery() FileQuery {
	return &fileQuery{d.db}
}

func (q *fileQuery) Create(UUID string, Name string, Size int64, Filename string) (*datastruct.File, error) {
	file := datastruct.File{
		UUID:     UUID,
		Name:     Name,
		Size:     Size,
		Filename: Filename,
	}
	err := q.db.Create(&file).Error
	return &file, err
}

func (q *fileQuery) Get(id uint) (*datastruct.File, error) {
	var file datastruct.File
	err := q.db.First(&file, id).Error
	return &file, err
}
