package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mada.h/educplus/models"
)

func connect() *gorm.DB {
	dsn := "root:20050412@tcp(127.0.0.1:3306)/educplus?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost);
	// mandrindra := User{
	//   Username: "educplus",
	//   Password: string(hash),
	// }
	// db.AutoMigrate(&models.User{}, &models.Event{}, &models.Mail{}, &models.Content{})
	// db.Create(&mandrindra);
	return db
}

func Authenticate(username, userPassword string) bool {
	userFound := GetUserByName(username)
	err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(userPassword))
	if err != nil {
		return false
	} else {
		return true
	}
}

func GetUserByName(username string) models.User {
	db := connect()
	var userFound models.User
	db.Table("users").Where("username = ?", username).First(&userFound)
	return userFound
}

func RegisterMail(mail models.Mail) {
	db := connect()
	db.Create(&mail)
}

func GetAllMail() []models.Mail {
	db := connect()
	emails := []models.Mail{}
	db.Find(&emails)
	return emails
}

func GetAllEvents() []models.Event {
	db := connect()
	events := []models.Event{}
	db.Table("events").Find(&events)
	return events
}

func RegisterEvent(event models.Event) {
	db := connect()
	db.Create(event)
}
