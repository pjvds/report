package command

import (
	"fmt"

	. "github.com/pjvds/slackme/context"
	"github.com/urfave/cli"
)

var List = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Action: func(c *cli.Context) error {
		context, err := LoadContext(c)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("failed to load context: %v", err), CONTEXT_ERR)
		}

		if context.NeedsLogin() {
			return ErrNeedLogin
		}

		fmt.Printf("name: %v\n", context.UserName)
		fmt.Printf("team: %v\n", context.UserName)
		fmt.Printf("email: %v\n", context.Email)

		if len(context.Channels) == 0 {
			fmt.Printf("\nno channels!\n")
		} else {
			fmt.Printf("\nchannels:\n")

			for _, channel := range context.Channels {
				fmt.Printf("\t%v\n", channel.Name)
			}
		}
		return nil
	},
}
