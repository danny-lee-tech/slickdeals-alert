package emailer

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

type GmailSettingConfig struct {
	SourceEmailAddress string
	TargetEmailAddress string
	Password           string
	Enabled            bool
}

func Email(config GmailSettingConfig, messageString string) error {
	if !config.Enabled {
		fmt.Println(config.Password)
		fmt.Println("Email notification disabled")
		return nil
	}

	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", config.SourceEmailAddress)
	message.SetHeader("To", config.TargetEmailAddress)
	message.SetHeader("Subject", "Slickdeals Alerts")

	// Set email body
	message.SetBody("text/plain", messageString)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, config.SourceEmailAddress, config.Password)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

	return nil
}
