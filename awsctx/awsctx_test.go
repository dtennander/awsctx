package awsctx

import (
	"gotest.tools/assert"
	"os"
	"strings"
	"testing"
)

var credFileContent []byte
var ctxFileContent []byte
var configFileContent []byte

var target *Awsctx

func before() {
	readFile = func(f string) (bs []byte, e error) {
		if strings.Contains(f, "credentials") {
			bs = []byte(`
				[default] 
				[OTHER_USER]`)
		} else if strings.Contains(f, "awsctx") {
			bs = []byte("USER")
		} else if strings.Contains(f, "config") {
			bs = []byte(`
				[default]
				[profile OTHER_USER]`)
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
		} else if strings.Contains(f, "config") {
			configFileContent = c
		}
		return nil
	}
	target, _ = New("aFolder")
}

func TestGetUsers(t *testing.T) {
	before()
	users, err := target.GetUsers()
	assert.NilError(t, err)
	assert.Equal(t, len(users), 2)
	assert.Equal(t, users[0].Name, "USER")
	assert.Equal(t, users[0].IsCurrent, true)
}

func TestSwitchUser(t *testing.T) {
	before()
	err := target.SwitchUser("OTHER_USER")
	assert.NilError(t, err)
	assert.Equal(t, string(credFileContent), `
				[USER] 
				[default]`)
	assert.Equal(t, string(configFileContent), `
				[profile USER]
				[default]`)
	assert.Equal(t, string(ctxFileContent), "OTHER_USER")
}

func TestRenameCtx(t *testing.T) {
	before()
	err := target.RenameUser("USER", "NEW_NAME")
	assert.NilError(t, err)
	assert.Equal(t, string(credFileContent), `
				[default] 
				[OTHER_USER]`)
	assert.Equal(t, string(configFileContent), `
				[default]
				[profile OTHER_USER]`)
	assert.Equal(t, string(ctxFileContent), "NEW_NAME")
}

func TestRenameNotCtx(t *testing.T) {
	before()
	err := target.RenameUser("OTHER_USER", "NEW_NAME")
	assert.NilError(t, err)
	assert.Equal(t, string(credFileContent), `
				[default] 
				[NEW_NAME]`)
	assert.Equal(t, string(configFileContent), `
				[default]
				[profile NEW_NAME]`)
}
