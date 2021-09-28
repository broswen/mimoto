package repository

import (
	"errors"
	"testing"
)

func TestMapRepository(t *testing.T) {
	mr, err := NewMap()
	if err != nil {
		t.Fatalf("NewMap: %v", err)
	}

	user := User{
		Email: "test@test.com",
		Name:  "test",
	}
	err = mr.Create(&user)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	err = mr.Create(&user)
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("Create: %v", err)
	}

	user, err = mr.FindByEmail(user.Email)
	if err != nil {
		t.Fatalf("FindByEmail: %v", err)
	}

	user, err = mr.FindByEmail("none")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("FindByEmail: %v", err)
	}

	err = mr.Save(&user)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

}
