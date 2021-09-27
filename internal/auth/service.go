package auth

import (
	"github.com/broswen/mimoto/internal/email"
	"github.com/broswen/mimoto/internal/repository"
)

type AuthService interface {
	Signup(email, password string) error
	Login(email, password string) error
	Confirm(email, code string) error
	SendReset(email string) error
	Reset(email, password, code string) error
}

type Service struct {
	emailService   email.EmailService
	userRepository repository.UserRepository
}

func New(emailService email.EmailService, userRepository repository.UserRepository) (Service, error) {
	return Service{
		emailService:   emailService,
		userRepository: userRepository,
	}, nil
}

func (s Service) Signup(email, password string) error {
	return nil
}

func (s Service) Login(email, password string) error {
	return nil
}

func (s Service) Confirm(email, code string) error {
	return nil
}

func (s Service) SendReset(email string) error {
	return nil
}

func (s Service) Reset(email, password, code string) error {
	return nil
}
