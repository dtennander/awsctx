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

func (c *configFile) renameProfile(oldName, newName string) error {
	var profileReg *regexp.Regexp
	if oldName == "default" {
		profileReg = defaultRegEx
	} else {
		var err error // needed to not override profileReg.
		profileReg, err = regexp.Compile(`\[profile ` + oldName + `]`)
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
	c.data = profileReg.ReplaceAll(c.data, []byte(newTag))
	return nil
}

var profilesRegex = regexp.MustCompile(`\[profile (\S+)]`)

func (c *configFile) getAllProfiles() []string {
	profiles := profilesRegex.FindAllSubmatch(c.data, -1)
	var result []string
	for i := range profiles {
		result = append(result, string(profiles[i][1]))
	}
	return result
}
