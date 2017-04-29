package main

import (
	"os"

	"github.com/pjvds/slackme/command"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := cli.App{
		Name: "slackme",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file, f",
				EnvVars: []string{"SLACKME_FILE"},
			},
			&cli.StringFlag{
				Name:    "host",
				EnvVars: []string{"SLACKME_HOST"},
				Value:   "https://slackme.org",
			},
		},
		Commands: []*cli.Command{
			command.Login,
			command.Add,
			command.Post,
			command.Exec,
			command.List,
		},
	}
	app.Run(os.Args)
}
