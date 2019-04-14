package awsctx

import "os"

type context struct {
	openFile
}

type NoContextError string

func (n NoContextError) Error() string {
	return string(n)
}

func newContext(folder string) (*context, error) {
	file, e := newOpenFile(folder + "/awsctx")
	if e != nil {
		if !os.IsNotExist(e) {
			return nil, e
		} else {
			return nil, NoContextError("awsctx file does not exist")
		}
	}
	return &context{openFile: *file}, nil
}

func createNewContext(folder, name string) error {
	return writeFile(folder + "/awsctx", []byte(name), 0644)
}

func (c *context) getContext() string {
	return string(c.data)
}

func (c *context) setContext(user string) {
	c.data = []byte(user)
}

func (c *context) isSet() bool {
	return len(c.data) != 0
}
