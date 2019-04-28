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

func (c *credentialsFile) getAllProfiles() []string {
	profiles := nameRegEx.FindAllSubmatch(c.data, -1)
	var result []string
	for i := range profiles {
		result = append(result, string(profiles[i][1]))
	}
	return result
}

func (c *credentialsFile) renameProfile(oldName, newName string) error {
	profileReg, err := regexp.Compile(`\[` + oldName + `]`)
	if err != nil {
		return err
	}
	c.data = profileReg.ReplaceAll(c.data, []byte("["+newName+"]"))
	return nil
}
