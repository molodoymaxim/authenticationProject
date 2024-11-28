package utils

import "fmt"

type EmailService struct {
	APIKey string
}

func NewEmailService(apiKey string) *EmailService {
	return &EmailService{
		APIKey: apiKey,
	}
}

func (e *EmailService) SendEmailWarning(userID, message string) {
	fmt.Printf("Sending email to user %s: %s\n", userID, message)
}
