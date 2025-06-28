package emailer

import (
	"encoding/base64"
	"fmt"
	"os"

	gomail "gopkg.in/mail.v2"
)

type Emailer struct {
	SMTP               string
	Port               int
	SourceEmailAddress string
	TargetEmailAddress string
	Subject            string
	PasswordFile       string
}

func (emailer Emailer) Email(messageString string) error {
	passwordBytes, err := os.ReadFile(emailer.PasswordFile)
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Email notification disabled")
		return nil
	}
	password, err := base64.StdEncoding.DecodeString(string(passwordBytes))
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Email notification disabled")
		return nil
	}

	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", emailer.SourceEmailAddress)
	message.SetHeader("To", emailer.TargetEmailAddress)
	message.SetHeader("Subject", emailer.Subject)

	// Set email body
	message.SetBody("text/plain", messageString)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(emailer.SMTP, emailer.Port, emailer.SourceEmailAddress, string(password))

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

	return nil
}
