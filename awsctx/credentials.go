package awsctx

import "regexp"

type credentials struct {
	openFile
}

func newCredentials(folder string) (*credentials, error) {
	file, e := newOpenFile(folder + "/credentials")
	return &credentials{openFile: *file}, e
}

func (c *credentials) getAllUsers() []string {
	nameRegEx := regexp.MustCompile(`\[(.+)]`)
	users := nameRegEx.FindAllSubmatch(c.data, -1)
	var result []string
	for i := range users {
		result = append(result, string(users[i][1]))
	}
	return result
}

func (c *credentials) renameUser(oldName, newName string) error {
	userReg, err := regexp.Compile(`\[` + oldName + `]`)
	if err != nil {
		return err
	}
	c.data = userReg.ReplaceAll(c.data, []byte("["+newName+"]"))
	return nil
}

func (c *credentials) userExists(user string) bool {
	userReg, err := regexp.Compile(`\[(` + user + `)]`)
	if err != nil {
		return false
	}
	return userReg.Match(c.data)
}
