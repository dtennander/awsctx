package main

import (
	"awsctx/awsctx"
	"log"
	"os"

	"github.com/urfave/cli"
)

var home = os.Getenv("HOME")

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		if len(c.Args()) > 0 {
			return awsctx.SwitchUser(home + "/.aws/credentials", c.Args()[0])
		} else {
			users, err := awsctx.GetUsers(home + "/.aws/credentials")
			if err != nil {
				return err
			}
			for i := range users  {
				print(users[i]+ "\n")
			}
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
