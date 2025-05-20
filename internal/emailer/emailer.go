package emailer

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func Email(messageString string) error {
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "purewhiteasian@gmail.com")
	message.SetHeader("To", "onfire_22043@yahoo.com")
	message.SetHeader("Subject", "Slickdeals Alerts")

	// Set email body
	message.SetBody("text/plain", messageString)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "purewhiteasian@gmail.com", "uvwd wtkw fprw ozkl")

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

	return nil
}
