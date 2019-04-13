package awsctx

import (
	"testing"
)

func TestGetUsers(t *testing.T) {
	readFile = func(f string) (bs []byte, e error) {
		bs = []byte(`[default] #name: DPT`)
		return
	}

	users, err := GetUsers("aFolder")
	if err != nil {
		t.Errorf("GetUser() error = %v", err)
	}

	if len(users) != 1 {
		t.Errorf("GetUser() should return one user got %d", len(users))
	}

	if users[0] != "USERNAME" {
		t.Errorf("GetUSer() did not return expected user USERNAME, got: %v", users[0])
	}

}