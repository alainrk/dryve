package datastruct

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Filename string
	UUID     string `gorm:"index:idx_uuid,unique"`
}
