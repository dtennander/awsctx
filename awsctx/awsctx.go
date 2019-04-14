package awsctx

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
)

type Context struct {
	Name      string
	IsCurrent bool
}

func GetUsers(folder string) ([]Context, error) {
	creds, err := newCredentials(folder)
	if err != nil {
		return nil, err
	}
	ctx, err := newContext(folder)
	if err != nil {
		return nil, err
	}
	users := creds.getAllUsers()
	var result []Context
	for _, user := range users {
		if user == "default" && ctx.isSet() {
			result = append(result, Context{Name: ctx.getContext(), IsCurrent: true})
		} else {
			result = append(result, Context{Name: user, IsCurrent: false})
		}
	}
	return result, nil
}

func SwitchUser(folder, user string) error {
	credentials, err := newCredentials(folder)
	if err != nil {
		return err
	}
	ctx, err := newContext(folder)
	if err != nil {
		return err
	}
	if !credentials.userExists(user) && user != ctx.getContext() {
		println("No user with the Name: \"" + user + "\".")
		return nil
	}
	if err := credentials.renameUser("default", ctx.getContext()); err != nil {
		return err
	}
	if err := credentials.renameUser(user, "default"); err != nil {
		return err
	}
	ctx.setContext(user)
	println("Switched to user \"" + user + "\".")
	if err =  credentials.store(); err != nil {
		return err
	}
	if err = ctx.store(); err != nil {
		return err
	}
	return nil
}

func RenameUser(folder, oldUser, newUser string) error {
	creds, err := newCredentials(folder)
	if err != nil {
		return err
	}
	ctx, err := newContext(folder)
	if err != nil {
		return err
	}
	switch {
	case oldUser == ctx.getContext():
		ctx.setContext(newUser)
	case creds.userExists(oldUser):
		if err := creds.renameUser(oldUser, newUser); err != nil {
			return err
		}
	default:
		println("No user with the Name: \"" + oldUser + "\".")
		return nil
	}
	println("Renamed user \"" + oldUser + "\" to \"" + newUser + "\".")
	if err =  creds.store(); err != nil {
		return err
	}
	if err = ctx.store(); err != nil {
		return err
	}
	return nil
}

var okName = regexp.MustCompile(`\S+`)

func SetUpDefaultContext(folder, defaultName string) error {
	if !okName.Match([]byte(defaultName)) {
		return errors.New(fmt.Sprintf("%s is not a valid context Name", defaultName))
	}
	return createNewContext(folder, defaultName)
}
