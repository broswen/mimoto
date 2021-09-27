package user

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/broswen/mimoto/internal/email"
	"github.com/broswen/mimoto/internal/repository"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo         repository.Repository
	emailService email.Service
}

func New(repo repository.Repository, emailService email.Service) (Service, error) {
	rand.Seed(time.Now().Unix())
	return Service{
		repo:         repo,
		emailService: emailService,
	}, nil
}

func (s Service) Signup(email, name, password string) error {
	// check if user exists, throw error if exists
	user := repository.User{Email: email, Name: name}
	tx := s.repo.DB.First(&user)

	if tx.Error == nil {
		return errors.New("user already exists with that email")
	}

	if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	code := generateCode()
	user.ConfirmationCode = code
	s.repo.DB.Create(&user)

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
	user := repository.User{Email: email}
	tx := s.repo.DB.First(&user)

	if tx.Error != nil {
		return tx.Error
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
	tx = s.repo.DB.Save(&user)

	if tx.Error != nil {
		return tx.Error
	}

	s.emailService.SendConfirmationSuccess(user.Name, user.Email)

	return nil
}

func (s Service) Login(email, password string) (string, string, error) {
	user := repository.User{Email: email}
	tx := s.repo.DB.First(&user)

	if tx.Error != nil {
		return "", "", tx.Error
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
	tx = s.repo.DB.Save(&user)

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
	// TODO check that the token actually matches
	// validate refresh token, throw error if invalid
	// find user from sub claim
	log.Println(email)
	log.Println(token)

	user := repository.User{Email: email}
	tx := s.repo.DB.First(&user)

	if tx.Error != nil {
		return "", tx.Error
	}

	if token != user.RefreshToken {
		return "", errors.New("refresh token doesn't match")
	}
	// throw error if not exist
	// check if refresh token matches, throw if not exist or doesn't match
	// generate new short lived jwt
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
	// validate token, throw error if invalid
	// find user from sub claim
	// throw error if not exist
	// set refresh token to "" on user in repo
	user := repository.User{Email: email}
	tx := s.repo.DB.First(&user)

	if tx.Error != nil {
		return tx.Error
	}

	user.RefreshToken = ""
	tx = s.repo.DB.Save(&user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s Service) SendReset(email string) error {
	// get account with email from sub
	user := repository.User{Email: email}
	tx := s.repo.DB.First(&user)

	if tx.Error != nil {
		return tx.Error
	}
	code := generateCode()
	// generate reset code
	// save code with user in repo
	user.ResetCode = code
	tx = s.repo.DB.Save(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return s.emailService.SendReset(user.Name, user.Email, code)
}

func (s Service) ResetPassword(email, password, code string) error {
	// get account with email from sub
	// throw error if not exist
	user := repository.User{Email: email}
	tx := s.repo.DB.First(&user)

	if tx.Error != nil {
		return tx.Error
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
	tx = s.repo.DB.Save(&user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
