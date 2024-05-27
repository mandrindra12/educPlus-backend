package mailsender

import (
  "github.com/go-mail/mail"
)

var api_key string  = "re_QGrwHe56_Prj2KQDpQX6ZDabuCZQeA4rK"

func SendMail(){
  m := mail.NewMessage();
  m.SetHeader("From", "mandrindraantonnio@gmail.com");
  m.SetHeader("To", "mandrindra1man@gmail.com");
  m.SetHeader("Subject", "Hello");
  m.SetBody("text/html", "hello <h1>Ao ve ?</h1>");
  //m.Attach(""); (joindre un fichier)
  d := mail.NewDialer("smtp.gmail.com", 587, "mandrindraantonnio@gmail.com", "whgjnecmwlzeybqf");
  
  if err := d.DialAndSend(m); err!=nil{
    println("Failed!!");
    panic(err);
  }
  println("Success!!");
}


