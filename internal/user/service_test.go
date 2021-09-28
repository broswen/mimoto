package user

import (
	"testing"

	"github.com/broswen/mimoto/internal/email"
	"github.com/broswen/mimoto/internal/repository"
)

func TestService(t *testing.T) {
	ur, _ := repository.NewMap()
	es, _ := email.NewConsole()
	us, err := New(ur, es)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	err = us.Signup("test@test.com", "test", "password")
	if err != nil {
		t.Fatalf("Signup: %v", err)
	}

	user, err := us.userRepository.FindByEmail("test@test.com")
	if err != nil {
		t.Fatalf("FindByEmail: %v", err)
	}
	if user.Email != "test@test.com" {
		t.Fatalf("email doesn't match: wanted %v but got %v", "test@test.com", user.Email)
	}

	if user.Confirmed {
		t.Fatalf("user is confirmed: wanted %v but got %v", false, user.Confirmed)
	}

	err = us.Confirm(user.Email, user.ConfirmationCode)
	if err != nil {
		t.Fatalf("Confirm: %v", err)
	}

	user, _ = ur.FindByEmail(user.Email)

	if !user.Confirmed {
		t.Fatalf("user is not confirmed: wanted %v but got %v", true, user.Confirmed)
	}

	token, refreshToken, err := us.Login(user.Email, "password")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if token == "" || refreshToken == "" {
		t.Fatalf("token or refreshToken is null")
	}

	err = us.SendReset(user.Email)
	if err != nil {
		t.Fatalf("SendReset: %v", err)
	}

	user, _ = ur.FindByEmail(user.Email)

	err = us.ResetPassword(user.Email, "newpassword", user.ResetCode)
	if err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}

}
