package models

import (
	"github.com/dgrijalva/jwt-go"
	_ "golang.org/x/text/date"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string
	jwt.StandardClaims
}

type Event struct {
	Title       string `json:"title"`
	ImagePath   string `json:"path"`
	Description string `json:"description"`
	// _ StartDate   date   `json:"start_date"`
	// _ EndDate     date   `json:"end_date"`
}

type CustomersMail struct {
	email string
}
