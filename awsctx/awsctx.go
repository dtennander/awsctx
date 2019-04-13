package awsctx

import (
	"bytes"
	"io/ioutil"
	"regexp"
)

var nameRegEx = regexp.MustCompile(`\[(.+)]`)
var currentUserRegEx = regexp.MustCompile(`\[(.+)](.*)#name:\s(.+)`)

var readFile = ioutil.ReadFile
var writeFile = ioutil.WriteFile

func GetUsers(folder string) ([]string, error) {
	content, err := readFile(folder + "/credentials")
	if err != nil {
		return nil, err
	}
	users := nameRegEx.FindAllSubmatch(content, -1)
	var result []string
	var defaultUserFound = false
	for i := range users {
		user := string(users[i][1])
		if user == "default" {
			defaultUserFound = true
			continue
		}
		result = append(result, string(user))
	}
	if currentUserRegEx.Match(content) {
		currentUser := string(currentUserRegEx.FindSubmatch(content)[3])
		result = append(result, currentUser)
	} else if defaultUserFound {
		result = append(result, "default")
	}
	return result, nil
}

func SwitchUser(folder, newUser string) error {
	file := folder + "/credentials"
	content, err := readFile(file)
	if err != nil {
		return err
	}
	newUserReg, err := regexp.Compile(`\[(` + newUser + `)]`)
	if err != nil {
		return err
	}
	if !newUserReg.Match(content) {
		print("No user with the name: \"" + newUser + "\".")
		return nil
	}
	content = currentUserRegEx.ReplaceAll(content, []byte("[$3]"))
	content = newUserReg.ReplaceAll(content, []byte("[default] #name: $1"))
	print("Switched to user \"" + newUser + "\".")
	return writeFile(file, content, 0644)
}

func RenameUser(folder, oldUser, newUser string) error {
	file := folder + "/credentials"
	content, err := readFile(file)
	if err != nil {
		return err
	}
	exp, err := regexp.Compile(`(\[|#name: )` + oldUser)
	if err != nil {
		return err
	}
	if !exp.Match(content) {
		print("No user with the name: \"" + oldUser + "\".")
		return nil
	}
	content = exp.ReplaceAll(content, []byte(`$1 REMOVE_ME` + newUser))
	content = bytes.Replace(content, []byte(" REMOVE_ME"), []byte(""), -1)
	print("Renamed user \"" + oldUser + "\" to \"" + newUser + "\".")
	return writeFile(file, content, 0644)
}
