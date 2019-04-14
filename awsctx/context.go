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
	file, e := newOpenFile(folder + "/awsctx")
	if e != nil {
		if !os.IsNotExist(e) {
			return nil, e
		}
		return nil, NoContextError("awsctx file does not exist")
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
