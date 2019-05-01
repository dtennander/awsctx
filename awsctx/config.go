package awsctx

import "regexp"

type configFile struct {
	openFile
}

func newConfigFile(folder string) (*configFile, error) {
	file, e := newOpenFile(folder + "/config")
	return &configFile{openFile: *file}, e
}

func (c *configFile) renameProfile(oldName, newName string) error {
	profileReg, err := getProfileRegex(oldName)
	if err != nil {
		return err
	}
	c.data = profileReg.ReplaceAll(c.data, getProfileTag(newName))

	sourceProfileRegEx, err := regexp.Compile(`source_profile *= *` + oldName)

	if err != nil {
		return err
	}
	c.data = sourceProfileRegEx.ReplaceAll(c.data, []byte("source_profile = "+newName))
	return nil
}

var defaultRegEx = regexp.MustCompile(`\[default]`)

func getProfileRegex(name string) (*regexp.Regexp, error) {
	if name == "default" {
		return defaultRegEx, nil
	} else {
		return regexp.Compile(`\[profile ` + name + `]`)
	}
}

func getProfileTag(name string) []byte {
	if name == "default" {
		return []byte("[default]")
	} else {
		return []byte("[profile " + name + "]")
	}
}

var profilesRegex = regexp.MustCompile(`(?:^|\n)\[(?:profile )?(\S+)]`)

func (c *configFile) getAllProfiles() []string {
	profiles := profilesRegex.FindAllSubmatch(c.data, -1)
	var result []string
	for i := range profiles {
		result = append(result, string(profiles[i][1]))
	}
	return result
}
