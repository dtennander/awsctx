package main

import (
	"awsctx/awsctx"
	"log"
	"os"

	"github.com/urfave/cli"
)

var awsFolder string

func main() {
	app := cli.NewApp()
	app.Name = "awsctx"
	app.Version = "0.1"
	app.HideVersion = true

	app.HelpName = "awsctx"
	app.ArgsUsage = "[user]"
	app.Usage = "A tool to switch aws user"
	app.UsageText = "awsctx [ <user> | rename <old user> <new user> ]"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "folder, f",
			Usage: "aws folder",
			Destination: &awsFolder,
			Value: os.Getenv("HOME") + "/.aws",
		},
	}

	app.Action = mainAction
	app.Commands = cli.Commands{
		{
			Name: "rename",
			ArgsUsage: "<old name> <new name>",
			Description: "renames a user to a new namer",
			ShortName: "r",
			Action: rename,

		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func rename(c *cli.Context) error {
	if c.NArg() != 2 {
		return cli.NewExitError("expected old name and new name", 1)
	}
	oldName := c.Args()[0]
	newName := c.Args()[1]
	return awsctx.RenameUser(awsFolder, oldName, newName)
}

func mainAction(c *cli.Context) error {
	switch c.NArg() {
	case 0:
		users, err := awsctx.GetUsers(awsFolder)
		if err != nil {
			return err
		}
		for i := range users {
			print(users[i] + "\n")
		}
		return nil
	case 1:
		return awsctx.SwitchUser(awsFolder, c.Args()[0])
	default:
		return cli.NewExitError("expected one or zero arguments", 1)
	}
}
