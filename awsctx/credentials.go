package awsctx

import "regexp"

type credentialsFile struct {
	openFile
}

func newCredentialsFile(folder string) (*credentialsFile, error) {
	file, e := newOpenFile(folder + "/credentials")
	return &credentialsFile{openFile: *file}, e
}

var nameRegEx = regexp.MustCompile(`\[(.+)]`)

func (c *credentialsFile) getAllUsers() []string {
	users := nameRegEx.FindAllSubmatch(c.data, -1)
	var result []string
	for i := range users {
		result = append(result, string(users[i][1]))
	}
	return result
}

func (c *credentialsFile) renameUser(oldName, newName string) error {
	userReg, err := regexp.Compile(`\[` + oldName + `]`)
	if err != nil {
		return err
	}
	c.data = userReg.ReplaceAll(c.data, []byte("["+newName+"]"))
	return nil
}

func (c *credentialsFile) userExists(user string) bool {
	userReg, err := regexp.Compile(`\[(` + user + `)]`)
	if err != nil {
		return false
	}
	return userReg.Match(c.data)
}
