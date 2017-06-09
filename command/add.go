package command

import (
	"fmt"

	. "github.com/pjvds/slackme/context"
	"gopkg.in/urfave/cli.v2"
)

var Add = &cli.Command{
	Name:  "add",
	Usage: "Add slackme to a Slack channel",
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
			message := fmt.Sprintf(
				"channel added succesfully, you can now post to %v:\n\n\t"+
					"slackme post '%v' 'hello world!'",
				channel.ChannelName, channel.ChannelName)
			return cli.Exit(message, 0)
		}
		return nil
	},
}
