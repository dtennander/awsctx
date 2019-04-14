package awsctx

type awsctx struct {
	credentialsFile *credentialsFile
	configFile      *configFile
	contextFile     *contextFile
}

func New(folder string) (*awsctx, error) {
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
	return &awsctx{
		credentialsFile: credFile,
		configFile:      configFile,
		contextFile:     contextFile,
	}, nil
}

type Context struct {
	Name      string
	IsCurrent bool
}

func (a *awsctx) GetUsers(folder string) ([]Context, error) {
	var result []Context
	for _, user := range a.credentialsFile.getAllUsers() {
		if user == "default" && a.contextFile.isSet() {
			result = append(result, Context{Name: a.contextFile.getContext(), IsCurrent: true})
		} else {
			result = append(result, Context{Name: user, IsCurrent: false})
		}
	}
	return result, nil
}

func (a *awsctx) SwitchUser(folder, user string) error {
	if !a.credentialsFile.userExists(user) && user != a.contextFile.getContext() {
		println("No user with the Name: \"" + user + "\".")
		return nil
	}
	if err := a.renameAll("default", a.contextFile.getContext()); err != nil {
		return err
	}
	if err := a.renameAll(user, "default"); err != nil {
		return err
	}
	a.contextFile.setContext(user)
	println("Switched to user \"" + user + "\".")
	return a.storeAll()
}

func (a *awsctx) renameAll(oldName, newName string) error {
	if err := a.credentialsFile.renameUser(oldName, newName); err != nil {
		return err
	}
	return a.configFile.renameUser(oldName, newName)
}

func (a *awsctx) storeAll() error {
	if err := a.credentialsFile.store(); err != nil {
		return err
	}
	if err := a.contextFile.store(); err != nil {
		return err
	}
	return a.configFile.store()
}

func (a *awsctx) RenameUser(folder, oldUser, newUser string) error {
	switch {
	case oldUser == a.contextFile.getContext():
		a.contextFile.setContext(newUser)
	case a.credentialsFile.userExists(oldUser):
		if err := a.renameAll(oldUser, newUser); err != nil {
			return err
		}
	default:
		println("No user with the Name: \"" + oldUser + "\".")
		return nil
	}
	println("Renamed user \"" + oldUser + "\" to \"" + newUser + "\".")
	return a.storeAll()
}

func SetUpDefaultContext(folder, defaultName string) error {
	return createNewContextFile(folder, defaultName)
}
