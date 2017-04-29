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
		cli.BoolFlag{
			Name: "code",
		},
	},
	BashComplete: func(c *cli.Context) {
		if c.NArg() == 0 {
			return
		}

		if context, err := LoadContext(c); err == nil {
			for _, channel := range context.Channels {
				println(channel.Name)
				return
			}
		}
	},
	Action: func(cli *cli.Context) error {
		if cli.NArg() != 2 {
			print("\"slack post\" requires 2 arguments.\nSee 'docker post --help'.\n\n" +
				"Usage:  docker post [OPTIONS] CHANNEL_NAME MESSAGE\n")
			return nil
		}

		channelName := cli.Args().First()
		if len(channelName) == 0 {
			print("\"slack post\" requires 2 arguments.\nSee 'docker post --help'.\n\n" +
				"Usage:  docker post [OPTIONS] CHANNEL_NAME MESSAGE\n")
			return nil
		}

		context, err := LoadContext(cli)
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

		message := cli.Args()[1]
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
