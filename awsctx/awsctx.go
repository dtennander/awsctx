package awsctx

import "github.com/DiTo04/awsctx/awsctx/strings"

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

func (a *Awsctx) GetUsers() ([]Context, error) {
	var result []Context
	for _, user := range strings.UnionOf(a.credentialsFile.getAllUsers(), a.configFile.getAllUsers()) {
		if user == "default" && a.contextFile.CurrentContext != "" {
			result = append(result, Context{Name: a.contextFile.CurrentContext, IsCurrent: true})
		} else {
			result = append(result, Context{Name: user, IsCurrent: false})
		}
	}
	return result, nil
}

func (a *Awsctx) RenameUser(oldUser, newUser string) error {
	switch {
	case oldUser == a.contextFile.CurrentContext:
		a.contextFile.CurrentContext = newUser
	case a.userExists(oldUser):
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

func (a *Awsctx) userExists(user string) bool {
	users := strings.UnionOf(a.configFile.getAllUsers(), a.credentialsFile.getAllUsers())
	return strings.Contains(users, user)
}

func (a *Awsctx) renameAll(oldName, newName string) error {
	if err := a.credentialsFile.renameUser(oldName, newName); err != nil {
		return err
	}
	return a.configFile.renameUser(oldName, newName)
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

func (a *Awsctx) SwitchUser(user string) error {
	if !a.userExists(user) && user != a.contextFile.CurrentContext {
		println("No user with the Name: \"" + user + "\".")
		return nil
	}
	if err := a.renameAll("default", a.contextFile.CurrentContext); err != nil {
		return err
	}
	if err := a.renameAll(user, "default"); err != nil {
		return err
	}
	a.contextFile.CurrentContext = user
	println("Switched to user \"" + user + "\".")
	return a.storeAll()
}

func SetUpDefaultContext(folder, defaultName string) error {
	return createNewContextFile(folder, defaultName)
}
