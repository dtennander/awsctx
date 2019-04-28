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

func setUpFiles(credentialProfiles []string, contextProfile *contextFile, configProfiles []string) {
	ctxBytes, _ := yaml.Marshal(contextProfile)
	readFile = func(f string) (bs []byte, e error) {
		if strings.Contains(f, "credentials") {
			bs = []byte(createFile("", credentialProfiles))
		} else if strings.Contains(f, "awsctx") {
			bs = ctxBytes
		} else if strings.Contains(f, "config") {
			bs = []byte(createFile("profile ", configProfiles))
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

func createFile(prefix string, profiles []string) string {
	var file strings.Builder
	for _, profile := range profiles {
		if profile == "default" {
			file.WriteString("[default]\n")
		} else {
			file.WriteString(fmt.Sprintf("[%s%s]\n", prefix, profile))
		}
	}
	return file.String()
}

func TestGetProfiles(t *testing.T) {
	profilesOnFile := []string{"default", "OTHER_PROFILE"}
	setUpFiles(profilesOnFile, &contextFile{CurrentProfile: "PROFILE"}, profilesOnFile)
	profiles, err := target.GetProfiles()
	assert.NilError(t, err)
	assert.Equal(t, len(profiles), 2)
	foundProfile := profiles[0].Name == "PROFILE" || profiles[1].Name == "PROFILE"
	assert.Assert(t, foundProfile)
}

func TestGetProfilesExtraInConfig(t *testing.T) {
	setUpFiles([]string{"default", "A"}, &contextFile{CurrentProfile: "PROFILE"}, []string{"default", "A", "B"})
	profiles, err := target.GetProfiles()
	assert.NilError(t, err)
	assert.Equal(t, len(profiles), 3)
}

func TestGetProfilesExtraInCredentials(t *testing.T) {
	setUpFiles([]string{"default", "A", "B"}, &contextFile{CurrentProfile: "PROFILE"}, []string{"default", "A"})
	profiles, err := target.GetProfiles()
	assert.NilError(t, err)
	assert.Equal(t, len(profiles), 3)
}

func TestSwitchProfile(t *testing.T) {
	profiles := []string{"default", "OTHER_PROFILE"}
	setUpFiles(profiles, &contextFile{CurrentProfile: "PROFILE"}, profiles)
	err := target.SwitchProfile("OTHER_PROFILE")
	assert.NilError(t, err)
	newProfiles := []string{"PROFILE", "default"}
	assert.Equal(t, string(credFileContent), createFile("", newProfiles))
	assert.Equal(t, string(configFileContent), createFile("profile ", newProfiles))
	newCtx := &contextFile{}
	_ = yaml.Unmarshal(ctxFileContent, newCtx)
	assert.Equal(t, newCtx.CurrentProfile, "OTHER_PROFILE")
	assert.Equal(t, newCtx.LastProfile, "PROFILE")
}

func TestSwitchBack(t *testing.T) {
	profiles := []string{"default", "OTHER_PROFILE"}
	setUpFiles(profiles, &contextFile{CurrentProfile: "PROFILE", LastProfile: "OTHER_PROFILE"}, profiles)
	err := target.SwitchBack()
	assert.NilError(t, err)
	newCtx := &contextFile{}
	_ = yaml.Unmarshal(ctxFileContent, newCtx)
	assert.Equal(t, newCtx.CurrentProfile, "OTHER_PROFILE")
	assert.Equal(t, newCtx.LastProfile, "PROFILE")
}

func TestRenameCtx(t *testing.T) {
	profiles := []string{"default", "OTHER_PROFILE"}
	setUpFiles(profiles, &contextFile{CurrentProfile: "PROFILE"}, profiles)
	err := target.RenameProfile("PROFILE", "NEW_NAME")
	assert.NilError(t, err)
	assert.Equal(t, string(credFileContent), createFile("", profiles))
	assert.Equal(t, string(configFileContent), createFile("profile ", profiles))
	assert.Equal(t, string(ctxFileContent), fmt.Sprintf("currentContext: %s\n", "NEW_NAME"))
}

func TestRenameNotCtx(t *testing.T) {
	profiles := []string{"default", "OTHER_PROFILE"}
	setUpFiles(profiles, &contextFile{CurrentProfile: "PROFILE"}, profiles)
	err := target.RenameProfile("OTHER_PROFILE", "NEW_NAME")
	assert.NilError(t, err)
	newProfiles := []string{"default", "NEW_NAME"}
	assert.Equal(t, string(credFileContent), createFile("", newProfiles))
	assert.Equal(t, string(configFileContent), createFile("profile ", newProfiles))
}
