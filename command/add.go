package command

import (
	"fmt"

	. "github.com/pjvds/slackme/context"
	"gopkg.in/urfave/cli.v2"
)

var Add = &cli.Command{
	Name: "add",
	Action: func(c *cli.Context) error {
		context, err := LoadContext(c)
		if err != nil {
			return cli.Exit(fmt.Sprintf("failed to load context: %v", err), CONTEXT_ERR)
		}

		channel, ok, err := context.AddChannel()
		if err != nil {
			return cli.Exit(fmt.Sprintf("failed to load context: %v", err), NEED_LOGIN)
		}
		if ok {
			fmt.Printf("channel added succesfully, run to post:\n\n\tslacke -c '%v' post", channel.ChannelName)
		}
		return nil
	},
}
