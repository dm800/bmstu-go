package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {
	// 'from' data
	from := "E"
	//from = "MAIL"
	password := "PASS"
	//password = "WORD"

	// 'to' email addresses
	to := []string{
		"EEEEMAIL",

		//You can add more than one 'to' email addresses
	}

	// smtp host server configuration
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"
	log.Println("Started log in")
	// Authentication.
	auth := smtp.PlainAuth(from, from, password, smtpHost)

	log.Println("Logged in")

	// Message.
	message :=
		"Subject: Оно должно работать \n" +
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
			"<br>" +
			"Обычный текст <br>" +
			"<p> Ваууу </p> "
	// Sending email.
	log.Println("Sending an email")
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	log.Println("Sent email")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}
