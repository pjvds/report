package command

import (
	"fmt"

	. "github.com/pjvds/slackme/context"
	"gopkg.in/urfave/cli.v2"
)

var List = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List all available Slack channels",
	Action: func(c *cli.Context) error {
		context, err := LoadContext(c)
		if err != nil {
			return cli.Exit(fmt.Sprintf("failed to load context: %v", err), CONTEXT_ERR)
		}

		if len(context.Channels) == 0 {
			fmt.Printf("no channels\n")
		} else {
			for _, channel := range context.Channels {
				fmt.Printf("%v/%v\n", channel.TeamName, channel.ChannelName)
			}
		}
		return nil
	},
}
