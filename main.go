package main

import (
	"awsctx/awsctx"
	"bufio"
	"github.com/urfave/cli"
	"log"
	"os"
)

const currentContextColor = "\033[1;33m"
const normalColor = "\033[0m"

var awsFolder string
var nameFlag string
var noColorFlag = false

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
			Name:        "folder, f",
			Usage:       "aws folder",
			Destination: &awsFolder,
			Value:       os.Getenv("HOME") + "/.aws",
		},
		cli.BoolFlag{
			Name: 	"no-color, nc",
			Usage: "remove color from output",
			Destination: &noColorFlag,
		},
	}

	app.Action = mainAction
	app.Commands = cli.Commands{
		{
			Name:        "rename",
			ArgsUsage:   "<old name> <new name>",
			Description: "renames a user to a new namer",
			ShortName:   "r",
			Action:      rename,
		},{
			Name: "setup",
			Description: "set up awsctx.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:"name, n",
					Usage: "set context name for default user",
					Destination: &nameFlag,
				},
			},
			Action: setup,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func rename(c *cli.Context) error {
	if c.NArg() != 2 {
		return cli.NewExitError("Expected old name and new name.", 1)
	}
	oldName := c.Args()[0]
	newName := c.Args()[1]
	return awsctx.RenameUser(awsFolder, oldName, newName)
}


func setup(c *cli.Context) error {
	var name string
	if nameFlag != "" {
		name = nameFlag
	} else {
		scanner := bufio.NewReader(os.Stdin)
		print("Name of current context: ")
		input, err := scanner.ReadString('\n')
		if err != nil {
			return err
		}
		name = input
	}
	err := awsctx.SetUpDefaultContext(awsFolder, name)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return err
}

func mainAction(c *cli.Context) error {
	switch c.NArg() {
	case 0:
		users, err := awsctx.GetUsers(awsFolder)
		if err != nil {
			_, ok := err.(awsctx.NoContextError)
			if !ok {
				return err
			}
			return cli.NewExitError("awsctx is not initialised. Please run `awsctx setup`.", 1)
		}
		for _, user := range users {
			var prefix string
			switch {
			case noColorFlag && user.IsCurrent:
				prefix = "*"
			case !noColorFlag && user.IsCurrent:
				prefix = currentContextColor
			case noColorFlag && !user.IsCurrent:
				prefix = " "
			case !noColorFlag && !user.IsCurrent:
				prefix = normalColor
			}
			print(prefix + user.Name + "\n")
		}
		return nil
	case 1:
		return awsctx.SwitchUser(awsFolder, c.Args()[0])
	default:
		return cli.NewExitError("expected one or zero arguments", 1)
	}
}
