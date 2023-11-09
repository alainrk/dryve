package repository

import (
	"dryve/internal/datastruct"

	"gorm.io/gorm"
)

type FileQuery interface {
	Create(UUID string, Name string, Size int64, Filename string) (*datastruct.File, error)
	Get(UUID string) (*datastruct.File, error)
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

func (q *fileQuery) Get(UUID string) (*datastruct.File, error) {
	var file datastruct.File

	err := q.db.Where("uuid = ?", UUID).First(&file).Error

	if err == gorm.ErrRecordNotFound {
		// TODO: Better position these errors
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}

	return &file, err
}
