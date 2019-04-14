package awsctx

import (
	"io/ioutil"
)

var readFile = ioutil.ReadFile
var writeFile = ioutil.WriteFile

type openFile struct {
	data     []byte
	fileName string
}

func newOpenFile(fileName string) (*openFile, error) {
	content, err := readFile(fileName)
	return &openFile{data: content, fileName: fileName}, err
}

func (oF *openFile) store() error {
	return writeFile(oF.fileName, oF.data, 0644)
}
