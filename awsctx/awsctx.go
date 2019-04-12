package awsctx

import (
	"errors"
	"io/ioutil"
	"regexp"
)

var nameRegEx = regexp.MustCompile(`\[(.+)]`)
var currentUserRegEx = regexp.MustCompile(`\[(.+)](.*)#name:\s(.+)`)


func GetUsers(file string) ([]string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	users := nameRegEx.FindAllSubmatch(content, -1)
	var result []string
	for i := range users {
		user := string(users[i][1])
		if user == "default" {
			continue
		}
		result = append(result, string(user))
	}
	currentUser := string(currentUserRegEx.FindSubmatch(content)[3])
	result = append(result, currentUser)
	return result, nil
}


func SwitchUser(file string, newUser string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	newUserReg, err := regexp.Compile(`\[(` + newUser + `)]`)
	if err != nil {
		return err
	}
	if newUserReg.Match(content) {
		content = currentUserRegEx.ReplaceAll(content, []byte("[$3]"))
		content = newUserReg.ReplaceAll(content, []byte("[default] #name: $1"))
		print("Switched to user \"" + newUser + "\".")
	} else {
		return errors.New("no user with the name: \"" + newUser + "\".")
	}
	return ioutil.WriteFile(file, content, 0644)
}
