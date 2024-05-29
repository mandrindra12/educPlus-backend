package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Claims struct {
	Username string
	Type     string
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Username string `gorm:"varchar(200)"`
	Password string `gorm:"varchar(200)"`
}

type Event struct {
	gorm.Model
	Title       string `gorm:"varchar(100)"`
	Description string `gorm:"varchar(500)"`
	EventType   string `gorm:"varchar(100)"`
	StartDate   time.Time
	EndDate     time.Time
	ImagePath   string `gorm:"varchar(200)"`
}

type Mail struct {
	gorm.Model
	EmailAddress string `json:"email"`
}

type Content struct {
	gorm.Model
	Path string
}
