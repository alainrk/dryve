package datastruct

import "gorm.io/gorm"

type File struct {
	gorm.Model
	// UUID of the file used for the filename
	UUID string `gorm:"index:idx_uuid,unique"`
	// Original filename
	Name string
	// Size of the file
	Size int64
	// Filename of the file on the server
	Filename string
}
