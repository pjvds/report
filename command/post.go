package command

import (
	"io/ioutil"
	"log"
	"os"

	. "github.com/pjvds/slackme/context"

	"github.com/urfave/cli"
)

var Post = cli.Command{
	Name: "post",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "c",
			EnvVar: "SLACKME_CHANNEL",
		},
		cli.BoolFlag{
			Name: "code",
		},
	},
	Action: func(cli *cli.Context) error {
		channelName := cli.String("c")
		if len(channelName) == 0 {
			log.Fatalf("missing channel name, please specify it like:\n\n\tslackme post -c '#general' 'hello from slackme!'\n\nOr set the SLACKME_CHANNEL environment variable.")
		}

		if len(cli.Args()) == 0 {
			log.Fatalf("missing message, please specify it like:\n\n\tslackme post -c '#general' 'hello from slackme!'")
		}

		if len(cli.Args()) > 1 {
			log.Fatalf("multiple arguments given, please specify your message as an single argument to post, like:\n\n\tslackme post -c '#general' 'hello from slackme!'")
		}

		context, err := LoadContext()
		if err != nil {
			log.Fatalf("failed to load context: %v", err)
		}

		if context.NeedsLogin() {
			log.Fatalf("not logged in, please run:\n\tslackme login")
		}

		channel, ok := context.ChannelByName(channelName)
		if !ok {
			log.Fatalf("channel not found, please run:\n\tslackme add '%v'", channelName)
		}

		message := cli.Args()[0]
		if message == "-" {
			stdin, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("failed to read from stdin: %v", err)
			}
			message = string(stdin)
		}

		if cli.Bool("code") {
			message = "```" + message + "```"
		}

		if err := channel.Post(message); err != nil {
			log.Fatalf("failed to post to channel %v: %v", channelName, err)
		}

		return nil
	},
}
