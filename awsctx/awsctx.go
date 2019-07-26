package awsctx

import "github.com/dtennander/awsctx/awsctx/strings"

type Awsctx struct {
	credentialsFile *credentialsFile
	configFile      *configFile
	contextFile     *contextFile
}

func New(folder string) (*Awsctx, error) {
	credFile, err := newCredentialsFile(folder)
	if err != nil {
		return nil, err
	}
	configFile, err := newConfigFile(folder)
	if err != nil {
		return nil, err
	}
	contextFile, err := newContextFile(folder)
	if err != nil {
		return nil, err
	}
	return &Awsctx{
		credentialsFile: credFile,
		configFile:      configFile,
		contextFile:     contextFile,
	}, nil
}

type Context struct {
	Name      string
	IsCurrent bool
}

func (a *Awsctx) GetProfiles() ([]Context, error) {
	var result []Context
	for _, profile := range strings.UnionOf(a.credentialsFile.getAllProfiles(), a.configFile.getAllProfiles()) {
		if profile == "default" && a.contextFile.CurrentProfile != "" {
			result = append(result, Context{Name: a.contextFile.CurrentProfile, IsCurrent: true})
		} else {
			result = append(result, Context{Name: profile, IsCurrent: false})
		}
	}
	return result, nil
}

func (a *Awsctx) RenameProfile(oldName, newName string) error {
	switch {
	case oldName == a.contextFile.CurrentProfile:
		a.contextFile.CurrentProfile = newName
	case a.profileExists(oldName):
		if err := a.renameAll(oldName, newName); err != nil {
			return err
		}
	default:
		println("No profile with the Name: \"" + oldName + "\".")
		return nil
	}
	println("Renamed profile \"" + oldName + "\" to \"" + newName + "\".")
	return a.storeAll()
}

func (a *Awsctx) profileExists(profile string) bool {
	profiles := strings.UnionOf(a.configFile.getAllProfiles(), a.credentialsFile.getAllProfiles())
	return strings.Contains(profiles, profile)
}

func (a *Awsctx) renameAll(oldName, newName string) error {
	if err := a.credentialsFile.renameProfile(oldName, newName); err != nil {
		return err
	}
	return a.configFile.renameProfile(oldName, newName)
}

func (a *Awsctx) storeAll() error {
	if err := a.credentialsFile.store(); err != nil {
		return err
	}
	if err := a.contextFile.store(); err != nil {
		return err
	}
	return a.configFile.store()
}

func (a *Awsctx) SwitchProfile(profile string) error {
	if !a.profileExists(profile) && profile != a.contextFile.CurrentProfile {
		println("No profile with the Name: \"" + profile + "\".")
		return nil
	}
	if err := a.renameAll("default", a.contextFile.CurrentProfile); err != nil {
		return err
	}
	if err := a.renameAll(profile, "default"); err != nil {
		return err
	}
	a.contextFile.LastProfile = a.contextFile.CurrentProfile
	a.contextFile.CurrentProfile = profile
	println("Switched to profile \"" + profile + "\".")
	return a.storeAll()
}

func (a *Awsctx) SwitchBack() error {
	return a.SwitchProfile(a.contextFile.LastProfile)
}

func SetUpDefaultContext(folder, defaultName string) error {
	return createNewContextFile(folder, defaultName)
}
