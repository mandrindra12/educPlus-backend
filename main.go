package main

import (
	"mada.h/educplus/server"
  // "mada.h/educplus/mailsender"
)

func main() {
  server.ListenAndServe()
  // mailsender.SendMail()
}
