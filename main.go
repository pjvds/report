package main

import (
	"os"

	"github.com/pjvds/slackme/command"
	"gopkg.in/urfave/cli.v2"
)

var version = "unknown"

func main() {

	app := cli.App{
		Name:                  "slackme",
		Usage:                 "post to slack from the commandline",
		Version:               version,
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file, f",
				EnvVars: []string{"SLACKME_FILE"},
				Value:   "$HOME/.slackme",
			},
			&cli.StringFlag{
				Name:    "host",
				EnvVars: []string{"SLACKME_HOST"},
				Value:   "https://slackme.org",
			},
		},
		Commands: []*cli.Command{
			command.Add,
			command.Post,
			command.Exec,
			command.List,
		},
	}
	app.Run(os.Args)
}
