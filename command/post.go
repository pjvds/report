package command

import (
	"bytes"
	"io"
	"log"
	"os"

	. "github.com/pjvds/slackme/context"

	"gopkg.in/urfave/cli.v2"
)

var Post = &cli.Command{
	Name:        "post",
	Description: "Post to a Slack channel.",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "code",
		},
	},
	ShellComplete: func(c *cli.Context) {
		if context, err := LoadContext(c); err == nil {
			for _, channel := range context.Channels {
				println(channel.ChannelName)
				return
			}
		}
	},
	Action: func(cli *cli.Context) error {
		if cli.NArg() != 2 {
			print("\"slackme post\" requires 2 arguments.\nSee 'docker post --help'.\n\n" +
				"Usage: slackme post [OPTIONS] CHANNEL_NAME MESSAGE\n")
			return nil
		}

		channelName := cli.Args().First()
		if len(channelName) == 0 {
			print("\"slackme post\" requires 2 arguments.\nSee 'docker post --help'.\n\n" +
				"Usage: slackme post [OPTIONS] CHANNEL_NAME MESSAGE\n")
			return nil
		}

		context, err := LoadContext(cli)
		if err != nil {
			log.Fatalf("failed to load context: %v", err)
		}

		channel, err := context.ChannelByName(channelName)
		if err != nil {
			return err
		}

		message := cli.Args().Get(1)

		if message == "-" {
			buffer := new(bytes.Buffer)

			_, err := io.Copy(io.MultiWriter(buffer, os.Stdout), os.Stdin)
			if err != nil {
				log.Fatalf("failed to read from stdin: %v", err)
			}
			message = buffer.String()
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
