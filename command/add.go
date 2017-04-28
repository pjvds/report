package command

import (
	"fmt"

	. "github.com/pjvds/slackme/context"
	"github.com/urfave/cli"
)

var Add = cli.Command{
	Name:    "add-channel",
	Aliases: []string{"ac"},
	Action: func(c *cli.Context) error {
		context, err := LoadContext(c)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("failed to load context: %v", err), CONTEXT_ERR)
		}

		if context.NeedsLogin() {
			return ErrNeedLogin
		}

		channel, ok, err := context.AddChannel()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("failed to load context: %v", err), NEED_LOGIN)
		}
		if ok {
			fmt.Printf("channel added succesfully, run to post:\n\n\tslacke -c '%v' post", channel.Name)
		}
		return nil
	},
}
