package user

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/broswen/mimoto/internal/email"
	"github.com/broswen/mimoto/internal/repository"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepository repository.UserRepository
	emailService   email.EmailService
}

func New(userRepository repository.UserRepository, emailService email.EmailService) (Service, error) {
	rand.Seed(time.Now().Unix())
	return Service{
		userRepository: userRepository,
		emailService:   emailService,
	}, nil
}

func (s Service) Signup(email, name, password string) error {
	user, err := s.userRepository.FindByEmail(email)
	if err == nil {
		return errors.New("user already exists with that email")
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		return err
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	code := generateCode()
	user.ConfirmationCode = code
	err = s.userRepository.Create(&user)
	if err != nil {
		return err
	}

	// send confirmation email
	err = s.emailService.SendConfirmation(name, email, code)
	if err != nil {
		return err
	}
	return nil
}

func generateCode() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func (s Service) Confirm(email, code string) error {
	// get user from repo, throw error if not exists
	user, err := s.userRepository.FindByEmail(email)

	if err != nil {
		return err
	}

	// check if confirmation code matches, throw error if not exist or not match
	if user.ConfirmationCode == "" {
		return errors.New("this user is already confirmed")
	}

	if code != user.ConfirmationCode {
		return errors.New("confirmation code doesn't match")
	}

	// set confirmed = true
	user.ConfirmationCode = ""
	user.Confirmed = true
	err = s.userRepository.Save(&user)

	if err != nil {
		return err
	}

	s.emailService.SendConfirmationSuccess(user.Name, user.Email)

	return nil
}

func (s Service) Login(email, password string) (string, string, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return "", "", err
	}

	if user.ConfirmationCode != "" || !user.Confirmed {
		return "", "", errors.New("account is not confirmed")
	}
	// hash password
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", "", err
	}

	tokenClaims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    "mimoto",
		Subject:   user.Email,
	}
	refreshTokenClaims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
		Issuer:    "mimoto",
		Subject:   user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", "", err
	}

	user.RefreshToken = signedRefreshToken
	err = s.userRepository.Save(&user)
	if err != nil {
		return "", "", err
	}
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func (s Service) Refresh(email, token string) (string, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if token != user.RefreshToken {
		return "", errors.New("refresh token doesn't match")
	}

	tokenClaims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    "mimoto",
		Subject:   user.Email,
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedToken, err := newToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s Service) Logout(email string) error {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return err
	}

	user.RefreshToken = ""
	err = s.userRepository.Save(&user)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) SendReset(email string) error {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return err
	}

	code := generateCode()
	user.ResetCode = code
	err = s.userRepository.Save(&user)
	if err != nil {
		return err
	}
	return s.emailService.SendReset(user.Name, user.Email, code)
}

func (s Service) ResetPassword(email, password, code string) error {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return err
	}

	if user.ResetCode == "" {
		return errors.New("no reset code for user")
	}

	if user.ResetCode != code {
		return errors.New("reset code doesn't match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// update user in repo with new password and set code to ""
	user.HashedPassword = string(hashedPassword)
	user.ResetCode = ""
	err = s.userRepository.Save(&user)
	if err != nil {
		return err
	}

	return nil
}
