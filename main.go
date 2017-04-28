package main

import (
	"os"

	"github.com/pjvds/slackme/command"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "slackme"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "file, f",
			EnvVar: "SLACKME_FILE",
		},
		cli.StringFlag{
			Name:   "host",
			EnvVar: "SLACKME_HOST",
			Value:  "https://slackme.org",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "channel",
			EnvVar: "SLACKME_CHANNEL",
		},
	}
	app.Commands = []cli.Command{
		command.Login,
		command.Add,
		command.Post,
		command.Exec,
		command.List,
	}
	app.Run(os.Args)
}
