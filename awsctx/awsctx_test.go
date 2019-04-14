package awsctx

import (
	"gotest.tools/assert"
	"os"
	"strings"
	"testing"
)

var credFileContent []byte
var ctxFileContent []byte

func initMocks() {
	readFile = func(f string) (bs []byte, e error) {
		if strings.Contains(f, "credentials") {
			bs = []byte(`
				[default] 
				[OTHER_USER]`)
		} else if strings.Contains(f, "awsctx") {
			bs = []byte("USER")
		}
		return
	}
	credFileContent = nil
	ctxFileContent = nil
	writeFile = func(f string, c []byte, p os.FileMode) error {
		if strings.Contains(f, "credentials") {
			credFileContent = c
		} else if strings.Contains(f, "awsctx") {
			ctxFileContent = c
		}
		return nil
	}
}

func TestGetUsers(t *testing.T) {
	initMocks()
	users, err := GetUsers("aFolder")
	assert.NilError(t, err)
	assert.Equal(t, len(users), 2)
	assert.Equal(t, users[0], "USER")
}

func TestSwitchUser(t *testing.T) {
	initMocks()
	err := SwitchUser("aFolder", "OTHER_USER")
	assert.NilError(t, err)
	assert.Equal(t, string(credFileContent), `
				[USER] 
				[default]`)
	assert.Equal(t, string(ctxFileContent), "OTHER_USER")
}

func TestRenameCtx(t *testing.T) {
	initMocks()
	err := RenameUser("aFolder", "USER", "NEW_NAME")
	assert.NilError(t, err)
	assert.Equal(t, string(credFileContent), `
				[default] 
				[OTHER_USER]`)
	assert.Equal(t, string(ctxFileContent), "NEW_NAME")
}

func TestRenameNotCtx(t *testing.T) {
	initMocks()
	err := RenameUser("aFolder", "OTHER_USER", "NEW_NAME")
	assert.NilError(t, err)
	assert.Equal(t, string(credFileContent), `
				[default] 
				[NEW_NAME]`)
}
