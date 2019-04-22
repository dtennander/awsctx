package awsctx

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"os"
	"strings"
	"testing"
)

var credFileContent []byte
var ctxFileContent []byte
var configFileContent []byte

var target *Awsctx

func setUpFiles(credentialUsers []string, contextUser *contextFile, configUsers []string) {
	ctxBytes, _ := yaml.Marshal(contextUser)
	readFile = func(f string) (bs []byte, e error) {
		if strings.Contains(f, "credentials") {
			bs = []byte(createFile("", credentialUsers))
		} else if strings.Contains(f, "awsctx") {
			bs = ctxBytes
		} else if strings.Contains(f, "config") {
			bs = []byte(createFile("profile ", configUsers))
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

func createFile(prefix string, users []string) string {
	var file strings.Builder
	for _, user := range users {
		if user == "default" {
			file.WriteString("[default]\n")
		} else {
			file.WriteString(fmt.Sprintf("[%s%s]\n", prefix, user))
		}
	}
	return file.String()
}

func TestGetUsers(t *testing.T) {
	usersOnFile := []string{"default", "OTHER_USER"}
	setUpFiles(usersOnFile, &contextFile{CurrentContext: "USER"}, usersOnFile)
	users, err := target.GetUsers()
	assert.NilError(t, err)
	assert.Equal(t, len(users), 2)
	foundUser := users[0].Name == "USER" || users[1].Name == "USER"
	assert.Assert(t, foundUser)
}

func TestGetUsersExtraInConfig(t *testing.T) {
	setUpFiles([]string{"default", "A"},&contextFile{CurrentContext: "USER"}, []string{"default", "A", "B"})
	users, err := target.GetUsers()
	assert.NilError(t, err)
	assert.Equal(t, len(users), 3)
}

func TestGetUsersExtraInCredentials(t *testing.T) {
	setUpFiles([]string{"default", "A", "B"},&contextFile{CurrentContext: "USER"}, []string{"default", "A"})
	users, err := target.GetUsers()
	assert.NilError(t, err)
	assert.Equal(t, len(users), 3)
}

func TestSwitchUser(t *testing.T) {
	users := []string{"default", "OTHER_USER"}
	setUpFiles(users, &contextFile{CurrentContext: "USER"}, users, )
	err := target.SwitchUser("OTHER_USER")
	assert.NilError(t, err)
	newUsers := []string{"USER", "default"}
	assert.Equal(t, string(credFileContent), createFile("", newUsers))
	assert.Equal(t, string(configFileContent), createFile("profile ", newUsers))
	newCtx := &contextFile{}
	_ = yaml.Unmarshal(ctxFileContent, newCtx)
	assert.Equal(t, newCtx.CurrentContext, "OTHER_USER")
	assert.Equal(t, newCtx.LastContext, "USER")
}

func TestSwitchBack(t *testing.T) {
	users := []string{"default", "OTHER_USER"}
	setUpFiles(users, &contextFile{CurrentContext: "USER", LastContext:"OTHER_USER"}, users, )
	err := target.SwitchBack()
	assert.NilError(t, err)
	newCtx := &contextFile{}
	_ = yaml.Unmarshal(ctxFileContent, newCtx)
	assert.Equal(t, newCtx.CurrentContext, "OTHER_USER")
	assert.Equal(t, newCtx.LastContext, "USER")
}

func TestRenameCtx(t *testing.T) {
	users := []string{"default", "OTHER_USER"}
	setUpFiles(users, &contextFile{CurrentContext: "USER"}, users, )
	err := target.RenameUser("USER", "NEW_NAME")
	assert.NilError(t, err)
	assert.Equal(t, string(credFileContent), createFile("", users))
	assert.Equal(t, string(configFileContent), createFile("profile ", users))
	assert.Equal(t, string(ctxFileContent), fmt.Sprintf("currentContext: %s\n", "NEW_NAME"))
}

func TestRenameNotCtx(t *testing.T) {
	users := []string{"default", "OTHER_USER"}
	setUpFiles(users, &contextFile{CurrentContext: "USER"}, users, )
	err := target.RenameUser("OTHER_USER", "NEW_NAME")
	assert.NilError(t, err)
	newUsers := []string{"default", "NEW_NAME"}
	assert.Equal(t, string(credFileContent), createFile("", newUsers))
	assert.Equal(t, string(configFileContent), createFile("profile ", newUsers))
}
