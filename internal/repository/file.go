package repository

import (
	"dryve/internal/datastruct"
	"time"

	"gorm.io/gorm"
)

type FileQuery interface {
	Create(UUID string, Name string, Size int64, Filename string) (datastruct.File, error)
	Get(UUID string) (datastruct.File, error)
	Delete(UUID string) error
	SearchByDateRange(from, to time.Time) ([]datastruct.File, error)
}

type fileQuery struct {
	db *gorm.DB
}

func (d *dao) NewFileQuery() FileQuery {
	return &fileQuery{d.db}
}

// Create a new file
func (q *fileQuery) Create(UUID string, Name string, Size int64, Filename string) (datastruct.File, error) {
	file := datastruct.File{
		UUID:     UUID,
		Name:     Name,
		Size:     Size,
		Filename: Filename,
	}
	err := q.db.Create(&file).Error
	return file, err
}

// Get a file by UUID
func (q *fileQuery) Get(UUID string) (datastruct.File, error) {
	var file datastruct.File
	err := q.db.Where("uuid = ?", UUID).First(&file).Error

	if err == gorm.ErrRecordNotFound {
		// TODO: Better position these errors
		return file, gorm.ErrRecordNotFound
	}
	if err != nil {
		return file, err
	}

	return file, err
}

// Search all files by date range
func (q *fileQuery) SearchByDateRange(from, to time.Time) ([]datastruct.File, error) {
	var files []datastruct.File

	err := q.db.Where("created_at BETWEEN ? AND ?", from, to).Find(&files).Error
	if err != nil {
		return nil, err
	}

	return files, err
}

// Delete a file by UUID
func (q *fileQuery) Delete(UUID string) error {
	err := q.db.Where("uuid = ?", UUID).Delete(&datastruct.File{}).Error
	return err
}
