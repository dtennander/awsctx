package awsctx

import "os"

type contextFile struct {
	openFile
}

type NoContextError string

func (n NoContextError) Error() string {
	return string(n)
}

func newContextFile(folder string) (*contextFile, error) {
	file, err := newOpenFile(folder + "/awsctx")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NoContextError("awsctx file does not exist")
		}
		return nil, err
	}
	return &contextFile{openFile: *file}, nil
}

func createNewContextFile(folder, name string) error {
	return writeFile(folder+"/awsctx", []byte(name), 0644)
}

func (c *contextFile) getContext() string {
	return string(c.data)
}

func (c *contextFile) setContext(user string) {
	c.data = []byte(user)
}

func (c *contextFile) isSet() bool {
	return len(c.data) != 0
}
