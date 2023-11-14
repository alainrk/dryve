package datastruct

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName   string
	LastName    string
	Email       string `gorm:"index:idx_email,unique"`
	Password    string
	PhoneNumber string
	Role        Role `gorm:"default:user"`
	Verified    bool
	EmailCode   string
}

type Role string

const (
	ADMIN Role = "admin"
	USER  Role = "user"
)
