package awsctx

import (
	"gopkg.in/yaml.v2"
	"os"
)

type contextFile struct {
	CurrentContext string `yaml:"currentContext"`
	LastContext string `yaml:"lastContext,omitempty"`
	filePath string
}

func (ctx *contextFile) store() error {
	out, err := yaml.Marshal(ctx)
	if err != nil {
		return err
	}
	return writeFile(ctx.filePath, out, 0644)
}

type NoContextError string

func (n NoContextError) Error() string {
	return string(n)
}

func newContextFile(folder string) (*contextFile, error) {
	filePath := folder + "/awsctx"
	file, err := readFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NoContextError("awsctx file does not exist")
		}
		return nil, err
	}
	ctx := &contextFile{}
	if err := yaml.Unmarshal(file, ctx); err != nil {
		return nil, err
	}
	ctx.filePath = filePath
	return ctx, nil
}

func createNewContextFile(folder, name string) error {
	newCtx := contextFile{CurrentContext:name}
	out, err := yaml.Marshal(newCtx)
	if err != nil {
		return err
	}
	return writeFile(folder+"/awsctx", out, 0644)
}
