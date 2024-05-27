package database;

import (
  "log"
  "time"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "golang.org/x/crypto/bcrypt"
)

type User struct {
  gorm.Model
  Username string `gorm:"varchar(200)"`
  Password string `gorm:"varchar(200)"`
}

type Event struct {
  gorm.Model
  Title string `gorm:"varchar(100)"`
  Description string `gorm:"varchar(500)"`
  EventType string `gorm:"varchar(100)"`
  StartDate time.Time
  EndDate time.Time
  ImagePath string `gorm:"varchar(200)"`
}

type Mail struct {
  gorm.Model
  EmailAddress string `gorm:"varchar(200)"`
}

type Content struct {
  gorm.Model
  Path string
}

func connect() *gorm.DB {
  dsn := "root:20050412@tcp(127.0.0.1:3306)/educplus?charset=utf8mb4&parseTime=True&loc=Local";
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{});
	if err != nil {
		log.Fatalf("failed to connect database: %v", err);
	}
  // hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost);
  // mandrindra := User{
  //   Username: "educplus",
  //   Password: string(hash),
  // }
  // db.AutoMigrate(&User{}, &Event{}, &Mail{}, &Content{});
  // db.Create(&mandrindra);
  return db;
}

func Authenticate(username, userPassword string) bool {
  userFound := GetUserByName(username);
  err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(userPassword));
  if err != nil {
    return false;
  } else {
    return true;
  }
}

func GetUserByName(username string) User {
  db := connect();
  var userFound User;
  db.Where("username = ?", username).First(&userFound);
  return userFound;
}

func RegisterMail(mail string) {
  db := connect();
  var email Mail;
  email = Mail{EmailAddress: mail};
  db.Create(&email);
}
