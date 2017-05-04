package command

import (
	"fmt"

	. "github.com/pjvds/slackme/context"
	"gopkg.in/urfave/cli.v2"
)

var Add = &cli.Command{
	Name:        "add",
	Description: "Add slackme to a Slack channel of your choice.",
	Action: func(c *cli.Context) error {
		context, err := LoadContext(c)
		if err != nil {
			return cli.Exit(fmt.Sprintf("failed to load context: %v\n", err), CONTEXT_ERR)
		}

		channel, ok, err := context.AddChannel()
		if err != nil {
			return err
		}
		if ok {
			fmt.Printf("channel added succesfully, run to post:\n\n\tslacke -c '%v' post", channel.ChannelName)
		}
		return nil
	},
}
