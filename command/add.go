package command

import (
	"fmt"
	"log"

	. "github.com/pjvds/slackme"
	"github.com/urfave/cli"
)

var Add = cli.Command{
	Name:    "add-channel",
	Aliases: []string{"ac"},
	Action: func(cli *cli.Context) error {
		context, err := LoadContext()
		if err != nil {
			log.Fatalf("failed to load context: %v", err)
		}

		if context.NeedsLogin() {
			log.Fatalf("need login, please run:\n\tslackme login")
		}

		channel, ok, err := context.AddChannel()
		if err != nil {
			log.Fatalf("failed to add channel: %v", err)
		}
		if ok {
			fmt.Printf("channel added succesfully, run to post:\n\tslacke -c '%v' post", channel.Name)
		}
		return nil
	},
}
