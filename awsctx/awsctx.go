package awsctx

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
)

func GetUsers(folder string) ([]string, error) {
	creds, err := newCredentials(folder)
	if err != nil {
		return nil, err
	}
	ctx, err := newContext(folder)
	if err != nil {
		return nil, err
	}
	users := creds.getAllUsers()
	var result []string
	for _, user := range users {
		if user == "default" && ctx.isSet() {
			user = ctx.getContext()
		}
		result = append(result, string(user))
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
		println("No user with the name: \"" + user + "\".")
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
		println("No user with the name: \"" + oldUser + "\".")
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
		return errors.New(fmt.Sprintf("%s is not a valid context name", defaultName))
	}
	return createNewContext(folder, defaultName)
}
