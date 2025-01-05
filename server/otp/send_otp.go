package otp

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendOTPEmail(recipientEmail, otp string) error {
	APIKey := os.Getenv("SENDGRID_API_KEY");
	if APIKey == "" {
		log.Println("No API KEY")
		return fmt.Errorf("no API KEY")
	}

	from := mail.NewEmail("Aethly", os.Getenv("MAIL_FROM"))
	subject := "Your Login OTP"
	plainTextContent := fmt.Sprintf("Your OTP is: %s. This OTP will expire in 5 minutes.", otp)
	htmlContent := fmt.Sprintf("<p>Your OTP is: <strong>%s</strong>. This OTP will expire in 5 minutes.</p>", otp)
	to := mail.NewEmail("Recipient", recipientEmail)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(APIKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	response, err := client.SendWithContext(ctx, message)
	if err != nil {
		return err
	}

	log.Printf("Email sent successfully. Status Code: %d", response.StatusCode)
	return nil
}
