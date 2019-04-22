package awsctx

import "regexp"

type configFile struct {
	openFile
}

func newConfigFile(folder string) (*configFile, error) {
	file, e := newOpenFile(folder + "/config")
	return &configFile{openFile: *file}, e
}

var defaultRegEx = regexp.MustCompile(`\[default]`)

func (c *configFile) renameUser(oldName, newName string) error {
	var userReg *regexp.Regexp
	if oldName == "default" {
		userReg = defaultRegEx
	} else {
		var err error // needed to not override userReg.
		userReg, err = regexp.Compile(`\[profile ` + oldName + `]`)
		if err != nil {
			return err
		}
	}
	var newTag string
	if newName == "default" {
		newTag = "[default]"
	} else {
		newTag = "[profile " + newName + "]"
	}
	c.data = userReg.ReplaceAll(c.data, []byte(newTag))
	return nil
}

var usersRegex = regexp.MustCompile(`\[profile (\S+)]`)

func (c *configFile) getAllUsers() []string {
	users := usersRegex.FindAllSubmatch(c.data, -1)
	var result []string
	for i := range users {
		result = append(result, string(users[i][1]))
	}
	return result
}
