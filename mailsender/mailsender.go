package mailsender

import (
	"github.com/go-mail/mail"
	"mada.h/educplus/models"
	"os"
)

var api_key string = "re_QGrwHe56_Prj2KQDpQX6ZDabuCZQeA4rK"

var password = string(os.Getenv("APP_PASSWORD"))

func SendMail(event models.Event, dest []models.Mail) error {
	m := mail.NewMessage()
	m.SetHeader("From", "mandrindraantonnio@gmail.com")
	// m.SetHeader("To", dest...)
	for _, email := range dest {
		m.SetHeader("To", email.EmailAddress)
	}
	m.SetHeader("Subject", event.Title)
	m.SetBody("text/plain", event.Description)
	// m.Attach(event.ImagePath) // (joindre un fichier)
	d := mail.NewDialer("smtp.gmail.com", 587, "mandrindraantonnio@gmail.com", "whgjnecmwlzeybqf")

	if err := d.DialAndSend(m); err != nil {
		return err
	} else {
		return nil
	}
}
