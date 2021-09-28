package email

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	Send(name, email, subject, text, html string) error
	SendConfirmation(name, email, code string) error
	SendConfirmationSuccess(name, email string) error
	SendReset(name, email, code string) error
}

type ConsoleService struct {
}

func NewConsole() (ConsoleService, error) {
	return ConsoleService{}, nil
}

func (cs ConsoleService) Send(name, email, subject, text, html string) error {
	from := fmt.Sprintf("noreply <%s>", os.Getenv("NOREPLY_EMAIL"))
	to := fmt.Sprintf("%s <%s>", name, email)
	fmt.Println(from, to, text)
	return nil
}

func (cs ConsoleService) SendConfirmation(name, email, code string) error {
	text := fmt.Sprintf("Please click this link to confirm your account.\n%s/confirm?email=%s&code=%s", os.Getenv("HOSTNAME"), email, code)

	return cs.Send(name, email, "Email Confirmation", text, text)
}

func (cs ConsoleService) SendConfirmationSuccess(name, email string) error {
	text := fmt.Sprintf("Your email was successfully confirmed!")

	return cs.Send(name, email, "Email Confirmation", text, text)
}

func (cs ConsoleService) SendReset(name, email, code string) error {
	text := fmt.Sprintf("Please click this link to reset your account password.\n%s/reset?email=%s&code=%s", os.Getenv("HOSTNAME"), email, code)

	return cs.Send(name, email, "Reset Password", text, text)
}

type SendGridService struct {
	sgClient *sendgrid.Client
}

func NewSendGrid() (SendGridService, error) {
	sgClient := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	return SendGridService{
		sgClient: sgClient,
	}, nil
}

func (s SendGridService) Send(name, email, subject, text, html string) error {
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

func (s SendGridService) SendConfirmation(name, email, code string) error {
	text := fmt.Sprintf("Please click this link to confirm your account.\n%s/confirm?email=%s&code=%s", os.Getenv("HOSTNAME"), email, code)

	return s.Send(name, email, "Email Confirmation", text, text)
}

func (s SendGridService) SendConfirmationSuccess(name, email string) error {
	text := fmt.Sprintf("Your email was successfully confirmed!")

	return s.Send(name, email, "Email Confirmation", text, text)
}

func (s SendGridService) SendReset(name, email, code string) error {
	text := fmt.Sprintf("Please click this link to reset your account password.\n%s/reset?email=%s&code=%s", os.Getenv("HOSTNAME"), email, code)

	return s.Send(name, email, "Reset Password", text, text)
}
