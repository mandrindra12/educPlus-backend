package mailsender;

import (
  "github.com/resend/resend-go/v2"
  "fmt"
)

func Resend() {

    client := resend.NewClient(api_key)

    params := &resend.SendEmailRequest{
        From:    "mandrindraantonnio@gmail.com",
        To:      []string{"mandrindraantonnio@gmail.com"},
        Html:    "<strong>hello world</strong>",
        Subject: "Hello from Golang",
        // Cc:      []string{"cc@example.com"},
        // Bcc:     []string{"bcc@example.com"},
        // ReplyTo: "replyto@example.com",
    }

    sent, err := client.Emails.Send(params)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println(sent.Id)
}

