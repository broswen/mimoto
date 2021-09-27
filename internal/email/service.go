package email

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Service struct {
	sgClient *sendgrid.Client
}

func New() (Service, error) {
	sgClient := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	return Service{
		sgClient: sgClient,
	}, nil
}

func (s Service) Send(name, email, subject, text, html string) error {
	from := mail.NewEmail("noreply", os.Getenv("NOREPLY_EMAIL"))
	to := mail.NewEmail(name, email)
	message := mail.NewSingleEmail(from, subject, to, text, html)
	response, err := s.sgClient.Send(message)
	log.Println(response)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) SendConfirmation(name, email, code string) error {
	text := fmt.Sprintf("Please click this link to confirm your account.\n%s/confirm?email=%s&code=%s", os.Getenv("HOSTNAME"), email, code)

	return s.Send(name, email, "Email Confirmation", text, text)
}

func (s Service) SendConfirmationSuccess(name, email string) error {
	text := fmt.Sprintf("Your email was successfully confirmed!")

	return s.Send(name, email, "Email Confirmation", text, text)
}

func (s Service) SendReset(name, email, code string) error {
	text := fmt.Sprintf("Please click this link to reset your account password.\n%s/reset?email=%s&code=%s", os.Getenv("HOSTNAME"), email, code)

	return s.Send(name, email, "Reset Password", text, text)
}
