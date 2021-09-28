package email

import "testing"

func TestConsoleService(t *testing.T) {
	cs, err := NewConsole()
	if err != nil {
		t.Fatalf("NewConsole: %v", err)
	}

	err = cs.SendConfirmation("test", "test@test.com", "12345")
	if err != nil {
		t.Fatalf("SendConfirmation: %v", err)
	}

	err = cs.SendConfirmationSuccess("test", "test@test.com")
	if err != nil {
		t.Fatalf("SendConfirmationSuccess: %v", err)
	}

	err = cs.SendReset("test", "test@test.com", "12345")
	if err != nil {
		t.Fatalf("SendReset: %v", err)
	}
}
